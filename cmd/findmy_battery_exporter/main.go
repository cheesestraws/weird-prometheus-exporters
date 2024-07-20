package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var prefix *string
var addr *string
var dump *bool

func produceMetricsBody() []byte {
	m := declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	devs, err := fetchDevices()
	if err != nil {
		log.Printf("err: %v", err)
	}
	ms := deviceMetricsFromDevices(devs)
	return m.Marshal(ms, nil)
}

func serve(addr string) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		w.Write(produceMetricsBody())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "metrics at /metrics")
	})

	log.Printf("listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	// flags
	prefix = flag.String("prefix", "findmy_", "prefix for metric names")
	addr = flag.String("addr", ":9405", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	flag.Parse()

	// check we're good to go
	warnForUntestedVersions()
	openFindMyApp()
	checkForAccess()

	if *dump {
		log.Printf("%s", produceMetricsBody())
	}

	serve(*addr)
}
