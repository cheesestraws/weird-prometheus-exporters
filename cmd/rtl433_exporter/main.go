package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var metrics Metrics

var baseURL *string
var prefix *string
var addr *string
var dump *bool
var useRFLookup *bool

func flushOldCrap() {
	for {
		time.Sleep(1 * time.Hour)
		metrics.Lock()
		metrics.DynamicMetrics.FlushOldCrap(7 * 24 * time.Hour)
		metrics.LastDynamicMetricFlush = time.Now().Unix()
		metrics.Unlock()
	}
}

func getBody() []byte {
	metrics.Lock()
	defer metrics.Unlock()

	m := declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	bs := m.Marshal(metrics, map[string]string{
		"base_url": *baseURL,
	})

	bs = append(bs, metrics.DynamicMetrics.PromBytes(*prefix+"evt_", *baseURL)...)

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
	baseURL = flag.String("baseurl", "http://rtl433.srv.lan:8433/", "base URL of rtl_433 instance")
	prefix = flag.String("prefix", "rtl_433_", "prefix for metric names")
	addr = flag.String("addr", ":9414", "address to listen on")
	useRFLookup = flag.Bool("rflookup", false, "use rflookup")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	flag.Parse()

	log.SetOutput(logWriter{})

	go fetchMetadata(*baseURL)
	go streamEvents(*baseURL)
	go flushOldCrap()

	serve(*addr)
}
