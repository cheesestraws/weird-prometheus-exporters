package main

import (
	"context"
	"flag"
	"log"
	"fmt"
	"time"
	"net/http"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var baseURL *string
var prefix *string
var addr *string
var dump *bool

func getBody() []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := fetchStatus(ctx, *baseURL)
	if err != nil {
		log.Printf("err: %v", err)
	}

	m := declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	return m.Marshal(s.Summarise(), map[string]string{
		"base_url": *baseURL,
	})
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
	baseURL = flag.String("baseurl", "http://127.0.0.1:9000", "base URL of Lyrion instance")
	prefix = flag.String("prefix", "lms_", "prefix for metric names")
	addr = flag.String("addr", ":9406", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	flag.Parse()

	if *dump {
		b := getBody()
		log.Printf("%s", b)
	}

	serve(*addr)
}
