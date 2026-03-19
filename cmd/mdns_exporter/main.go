package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
	_ "github.com/cheesestraws/weird-prometheus-exporters/lib/logutil"
)

var watcher *mdnsWatcher

var prefix *string
var addr *string
var dump *bool

func getBody() []byte {
	m := declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	bs := m.Marshal(watcher.metrics(), nil)

	return bs
}

func serve(addr string) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		bs := getBody()
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

func main() {
	prefix = flag.String("prefix", "mdns_", "prefix for metric names")
	addr = flag.String("addr", ":9415", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	flag.Parse()

	watcher = newmdnsWatcher()
	go watcher.watchServices()

	serve(*addr)
}
