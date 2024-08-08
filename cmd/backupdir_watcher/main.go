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

	"github.com/cheesestraws/scb/lib/metadata"
	
	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var pathParams multiflag = multiflag{}
var scbdir *string
var scbname string
var dirs map[string]BackupPath

type BackupPath struct {
	Path string
	SCB bool
}

func makePathMap(ps []string) map[string]BackupPath {
	pm := make(map[string]BackupPath)

	for _, p := range ps {
		path := filepath.Clean(p)
		basename := filepath.Base(p)

		pm[basename] = BackupPath{
			Path: path,
		}
	}

	return pm
}

type mostRecent struct {
	path             string
	partiallyErrored bool

	foundNewest bool
	when        time.Time
	filename    string
	size        int64

	foundPrevious    bool
	previousWhen     time.Time
	previousFilename string
	previousSize     int64
}

var re *regexp.Regexp = regexp.MustCompile(`\d{8}`)

func findMostRecent(path BackupPath) (mostRecent, error) {
	var m mostRecent
	m.path = path.Path

	entries, err := os.ReadDir(path.Path)
	if err != nil {
		return m, err
	}

	// files are sorted by filename, so iterate *backwards*
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		
		if strings.HasSuffix(e.Name(), ".scbmeta") {
			continue
		}

		date := re.FindString(e.Name())
		if date != "" && !m.foundNewest {
			m.when, _ = time.Parse("20060102", date)
			m.filename = e.Name()
			m.foundNewest = true

			s, err := os.Stat(filepath.Join(path.Path, e.Name()))
			if err != nil {
				m.partiallyErrored = true
				log.Printf("err: %v", err)
			}
			m.size = s.Size()

			continue
		} else if date != "" {
			m.previousWhen, _ = time.Parse("20060102", date)
			m.previousFilename = e.Name()
			m.foundPrevious = true

			s, err := os.Stat(filepath.Join(path.Path, e.Name()))
			if err != nil {
				m.partiallyErrored = true
				log.Printf("err: %v", err)
			}
			m.previousSize = s.Size()

			break
		}
	}

	if !m.foundNewest {
		return m, errors.New("no backup found")
	} else {
		return m, nil
	}
}

type scbmeta struct {
	Backup string `prometheus_label:"backup"`
	BackupHost string `prometheus_label:"scb_host"`
	Kind string `prometheus_label:"scb_kind"`
}

type BackupAges struct {
	Timestamp    map[string]int     `prometheus_map:"timestamp" prometheus_map_key:"backup" prometheus_help:"UNIX timestamp"`
	Age          map[string]int     `prometheus_map:"age" prometheus_map_key:"backup" prometheus_help:"in seconds"`
	Interval     map[string]int     `prometheus_map:"last_interval" prometheus_map_key:"backup" prometheus_help:"in seconds"`
	Size         map[string]int64   `prometheus_map:"size" prometheus_map_key:"backup" prometheus_help:"in bytes"`
	SizeDelta    map[string]int64   `prometheus_map:"size_change" prometheus_map_key:"backup" prometheus_help:"in bytes"`
	SizeDeltaPct map[string]float64 `prometheus_map:"size_change_pct" prometheus_map_key:"backup" prometheus_help:"percentage"`
	Error        map[string]int     `prometheus_map:"errors" prometheus_map_key:"backup"`
	
	Active map[string]int `prometheus_map:"active" prometheus_map_key:"backup"`

	SCBOK     int `prometheus:"scb_lastrun_ok"`
	SCBErrors int `prometheus:"scb_lastrun_errors"`
	
	SCBFetchTime map[string]float64 `prometheus_map:"scb_fetch_time" prometheus_map_key:"backup" prometheus_help:"in seconds"`
	SCBMeta map[scbmeta]int `prometheus_map:"scb_meta"`
}

func gatherBackupAges(paths map[string]BackupPath) BackupAges {
	ages := BackupAges{
		Timestamp:    make(map[string]int),
		Age:          make(map[string]int),
		Interval:     make(map[string]int),
		Size:         make(map[string]int64),
		SizeDelta:    make(map[string]int64),
		SizeDeltaPct: make(map[string]float64),
		Error:        make(map[string]int),
		
		SCBFetchTime: make(map[string]float64),
		SCBMeta: make(map[scbmeta]int),
	}

	for k, v := range paths {
		m, err := findMostRecent(v)
		if err != nil {
			ages.Error[k] = 1
			log.Printf("%v: %v", k, err)
			continue
		}

		if m.partiallyErrored {
			ages.Error[k] = 1
		} else {
			ages.Error[k] = 0
		}

		ages.Timestamp[k] = int(m.when.Unix())
		ages.Age[k] = int(time.Now().Sub(m.when).Seconds())
		if m.foundPrevious {
			ages.Interval[k] = int(m.when.Sub(m.previousWhen).Seconds())
		}

		ages.Size[k] = m.size

		if m.foundPrevious {
			ages.SizeDelta[k] = m.size - m.previousSize
			if m.previousSize == 0 {
				m.previousSize = 1
			}
			ages.SizeDeltaPct[k] = (float64(m.size-m.previousSize) / float64(m.previousSize)) * 100
		}
		
		if v.SCB {
			dealWithSCBMetadata(k, filepath.Join(v.Path, m.filename), &ages)
		}
	}

	dealWithSCBLogs(&ages)

	return ages
}

func gatherOurBackupMetadata(paths map[string]BackupPath, into *BackupAges) {
	into.Active = make(map[string]int)
	
	for k, v := range paths {
		inactiveFile := filepath.Join(v.Path, ".inactive")
		_, err := os.Stat(inactiveFile)
		
		if err != nil {
			into.Active[k] = 1
		} else {
			into.Active[k] = 0
		}
	}
}

func dealWithSCBMetadata(backup string, latest string, into *BackupAges) error {
	md, err := metadata.ReadFor(latest)
	if err != nil {
		return err
	}
	if md == nil {
		return nil
	}
	
	m := scbmeta{
		Backup: backup,
		BackupHost: md.BackupHost,
		Kind: md.Kind,
	}
	
	into.SCBMeta[m] = 1
	into.SCBFetchTime[backup] = md.FetchTime.Seconds()
	return nil
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

func discoverSCBBackups(into *map[string]BackupPath) {
	if *scbdir == "" {
		return
	}

	entries, err := os.ReadDir(*scbdir)
	if err != nil {
		log.Printf("scb discovery error: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			(*into)[filepath.Join(scbname, e.Name())] = BackupPath{
				Path: filepath.Join(*scbdir, e.Name()),
				SCB: true,
			}
		}
	}
}

func makeBody() []byte {
	checkPaths := maps.Clone(dirs)
	discoverSCBBackups(&checkPaths)

	ba := gatherBackupAges(checkPaths)
	gatherOurBackupMetadata(checkPaths, &ba)

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

	dirs = makePathMap([]string(pathParams))
	if *scbdir != "" {
		*scbdir = filepath.Clean(*scbdir)
		scbname = filepath.Base(*scbdir)
	}

	if *dump {
		log.Printf("%s", makeBody())
	}

	serve(*addr)
}
