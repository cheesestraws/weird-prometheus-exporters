package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var addr *string
var endpoint *string
var prefix *string
var display *string
var marshaller declprom.Marshaller

type xsetStatus struct {
	DisplayOn int `prometheus:"monitor_on"`
}

func getstatus() xsetStatus {
	stats := xsetStatus{
		DisplayOn: -1,
	}
	
	cmd := exec.Command("xset", "q")
	cmd.Env = append(os.Environ(), "DISPLAY=" + *display)
	
	var out []byte
	var err error
	if out, err = cmd.Output(); err != nil {
		log.Printf("error running xset: %v", err)
		return stats
	}
	
	switch {
	case bytes.Contains(out, []byte("Monitor is On")):
		stats.DisplayOn = 1
	case bytes.Contains(out, []byte("Monitor is Off")):
		stats.DisplayOn = 0
	}
	
	return stats
}

func serve(addr string) {
	http.HandleFunc(*endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		w.Write(marshaller.Marshal(getstatus(), nil))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "ok")
	})

	log.Printf("listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}


func main() {
	addr = flag.String("addr", ":9413", "address to listen on")
	endpoint = flag.String("endpoint", "/metrics", "the metrics endpoint")
	prefix = flag.String("prefix", "xset_", "prefix for metric names")
	display = flag.String("display", ":0", "X display to use")
	flag.Parse()
	
	marshaller = declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	serve(*addr)
}
