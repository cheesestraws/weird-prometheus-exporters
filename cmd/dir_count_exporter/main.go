package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var addr *string
var dump *bool
var endpoint *string
var prefix *string
var globs []string
var marshaller declprom.Marshaller

func serve(addr string) {
	http.HandleFunc(*endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		matches, errs := genGlobs(globs)
		stats := gatherStats(matches, errs)

		w.Write(marshaller.Marshal(stats, nil))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "ok")
	})

	log.Printf("listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	addr = flag.String("addr", ":9411", "address to listen on")
	endpoint = flag.String("endpoint", "/metrics", "the metrics endpoint")
	prefix = flag.String("prefix", "dir_", "prefix for metric names")
	flag.Parse()
	
	if len(flag.Args()) < 1 {
		log.Printf("no glob, pls provide")
		return
	}
	
	globs = flag.Args()	
	
	marshaller = declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	serve(*addr)
}
