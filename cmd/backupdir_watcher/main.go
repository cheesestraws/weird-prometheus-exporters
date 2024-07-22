package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"net/http"
	"fmt"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var pathParams multiflag = multiflag{}
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
			log.Printf("%s: %v", path, t)

			return t, nil
		}
	}

	return time.Time{}, errors.New("no backup found")
}

type BackupAges struct {
	Timestamp map[string]int `prometheus_map:"timestamp" prometheus_map_key:"backup" prometheus_help:"UNIX timestamp"`
	Age       map[string]int `prometheus_map:"age" prometheus_map_key:"backup" prometheus_help:"in seconds"`
	Error     map[string]int `prometheus_map:"errors" prometheus_map_key:"backup"`
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
		}

		ages.Timestamp[k] = int(t.Unix())
		ages.Age[k] = int(time.Now().Sub(t).Seconds())
	}

	return ages
}

func makeBody() []byte {
	ba := gatherBackupAges(dirs)
	
	log.Printf("%+v", ba)
	
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
	prefix = flag.String("prefix", "backupdir_", "prefix for metric names")
	addr = flag.String("addr", ":9407", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")

	flag.Parse()

	dirs = cleanPaths([]string(pathParams))

	if *dump {
		log.Printf("%s", makeBody())
	}
	
	serve(*addr)
}
