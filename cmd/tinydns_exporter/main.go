package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var datafile *string
var datacdb *string
var servicedir *string
var logdir *string
var suffix *string
var verbose *bool
var addr *string
var dump *bool
var endpoint *string

func serve(addr string) {
	http.HandleFunc(*endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		body := getBody()

		if *dump {
			log.Printf("%s", body)
		}

		w.Write(body)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "ok")
	})

	log.Printf("listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	addr = flag.String("addr", ":9408", "address to listen on")
	datafile = flag.String("data", "", "path to the text data file")
	datacdb = flag.String("datacdb", "", "path to the data.cdb file")
	servicedir = flag.String("servicedir", "", "path to the tinydns service dir")
	logdir = flag.String("logdir", "", "path to the current tinydns logs")
	suffix = flag.String("suffix", "", "a common suffix for our dns server names")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	verbose = flag.Bool("v", false, "verbose mode")
	endpoint = flag.String("endpoint", "/metrics", "the metrics endpoint")

	flag.Parse()

	if *datafile == "" && *datacdb == "" && *servicedir == "" && *logdir == "" {
		log.Printf("you probably want to supply at least one thing to watch")
		return
	}

	start()
	serve(*addr)
}
