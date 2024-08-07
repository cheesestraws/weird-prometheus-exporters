package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
	"github.com/cheesestraws/weird-prometheus-exporters/lib/sas2ircu"
)

var addr *string
var dump *bool
var endpoint *string
var executable *string
var prefix *string

var state struct {
	l  sync.RWMutex
	bs []byte
}

func gensummaries() *Summary {
	q := sas2ircu.Querier{
		Executable: *executable,
	}

	details, err := q.Get()
	if err != nil {
		log.Printf("err: %v", err)
		return &Summary{FetchError: 1}
	}

	return Summarise(details)
}

func pollLoop() {
	for {
		sums := gensummaries()

		m := declprom.Marshaller{
			MetricNamePrefix: *prefix,
		}

		bs := m.Marshal(*sums, nil)

		if *dump {
			log.Printf("%s", bs)
		}

		state.l.Lock()
		state.bs = bs
		state.l.Unlock()

		time.Sleep(5 * time.Minute)
	}
}

func serve(addr string) {
	http.HandleFunc(*endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		state.l.RLock()
		body := state.bs
		state.l.RUnlock()

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
	addr = flag.String("addr", ":9409", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	endpoint = flag.String("endpoint", "/metrics", "the metrics endpoint")
	executable = flag.String("executable", "/usr/local/bin/sas2ircu", "path to the sas2ircu executable")
	prefix = flag.String("prefix", "sas2ircu_", "prefix for metric names")
	flag.Parse()

	go pollLoop()
	serve(*addr)
}
