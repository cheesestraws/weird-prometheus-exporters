package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var pathParams multiflag = multiflag{}
var scbdir *string
var scbname string
var dirs map[string]string

func cleanPaths(ps []string) map[string]string {
	pm := make(map[string]string)

	for _, p := range ps {
		path := filepath.Clean(p)
		basename := filepath.Base(p)

		pm[basename] = path
	}

	return pm
}

var re *regexp.Regexp = regexp.MustCompile(`\d{8}`)

func findMostRecent(path string) (time.Time, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return time.Time{}, err
	}

	// files are sorted by filename, so iterate *backwards*
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]

		date := re.FindString(e.Name())
		if date != "" {
			t, _ := time.Parse("20060102", date)
			return t, nil
		}
	}

	return time.Time{}, errors.New("no backup found")
}

type BackupAges struct {
	Timestamp map[string]int `prometheus_map:"timestamp" prometheus_map_key:"backup" prometheus_help:"UNIX timestamp"`
	Age       map[string]int `prometheus_map:"age" prometheus_map_key:"backup" prometheus_help:"in seconds"`
	Error     map[string]int `prometheus_map:"errors" prometheus_map_key:"backup"`

	SCBOK     int `prometheus:"scb_lastrun_ok"`
	SCBErrors int `prometheus:"scb_lastrun_errors"`
}

func gatherBackupAges(paths map[string]string) BackupAges {
	ages := BackupAges{
		Timestamp: make(map[string]int),
		Age:       make(map[string]int),
		Error:     make(map[string]int),
	}

	for k, v := range paths {
		t, err := findMostRecent(v)
		if err != nil {
			ages.Error[k] = 1
			log.Printf("%v: %v", k, err)
		} else {
			ages.Error[k] = 0
		}

		ages.Timestamp[k] = int(t.Unix())
		ages.Age[k] = int(time.Now().Sub(t).Seconds())
	}

	dealWithSCBLogs(&ages)

	return ages
}

func dealWithSCBLogs(into *BackupAges) error {
	seenOK := false
	errors := 0

	logfile, err := findRecentSCBLog()
	if err != nil {
		return err
	}

	bs, err := os.ReadFile(filepath.Join(*scbdir, "log", logfile))
	ls := strings.Split(string(bs), "\n")
	for _, l := range ls {
		if l == "ok" {
			seenOK = true
		} else if l != "" {
			errors++
		}
	}

	if seenOK {
		into.SCBOK = 1
	} else {
		into.SCBOK = 0
	}
	into.SCBErrors = errors

	return nil
}

func findRecentSCBLog() (string, error) {
	entries, err := os.ReadDir(filepath.Join(*scbdir, "log"))
	if err != nil {
		return "", err
	}

	// files are sorted by filename, so iterate *backwards*
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]

		date := re.FindString(e.Name())
		if date != "" {
			return e.Name(), nil
		}
	}

	return "", errors.New("no log found")
}

func discoverSCBBackups(into *map[string]string) {
	if *scbdir == "" {
		return
	}

	entries, err := os.ReadDir(*scbdir)
	if err != nil {
		log.Printf("scb discovery error: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			(*into)[filepath.Join(scbname, e.Name())] = filepath.Join(*scbdir, e.Name())
		}
	}
}

func makeBody() []byte {
	checkPaths := maps.Clone(dirs)
	discoverSCBBackups(&checkPaths)

	ba := gatherBackupAges(checkPaths)

	m := declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	return m.Marshal(ba, nil)
}

func serve(addr string) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		bs := makeBody()
		if *dump {
			log.Printf("%s", bs)
		}

		w.Write(bs)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "metrics at /metrics")
	})

	log.Printf("listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

var prefix *string
var addr *string
var dump *bool

func main() {
	flag.Var(&pathParams, "dir", "directory/ies to watch (use param more than once")
	scbdir = flag.String("scbdir", "", "scb -into directory to watch")
	prefix = flag.String("prefix", "backupdir_", "prefix for metric names")
	addr = flag.String("addr", ":9407", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")

	flag.Parse()

	dirs = cleanPaths([]string(pathParams))
	if *scbdir != "" {
		*scbdir = filepath.Clean(*scbdir)
		scbname = filepath.Base(*scbdir)
	}

	if *dump {
		log.Printf("%s", makeBody())
	}

	serve(*addr)
}
