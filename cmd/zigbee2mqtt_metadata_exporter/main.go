package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	
	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var addr *string
var endpoint *string
var prefix *string
var mapNHrs *int
var basetopic *string
var marshaller declprom.Marshaller

func serve(addr string) {
	http.HandleFunc(*endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		
		promstate.Lock()
		defer promstate.Unlock()
		w.Write(marshaller.Marshal(promstate, map[string]string{
			"base_topic": *basetopic,
		}))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "ok")
	})

	log.Printf("listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}


func main() {
	broker := flag.String("broker", "", "mqtt broker to connect to")
	user := flag.String("user", "", "username")
	pass := flag.String("pass", "", "password")
	basetopic = flag.String("basetopic", "", "base topic")
	addr = flag.String("addr", ":9412", "address to listen on")
	endpoint = flag.String("endpoint", "/metrics", "the metrics endpoint")
	prefix = flag.String("prefix", "zigbee_", "prefix for metric names")
	mapNHrs = flag.Int("map_every", 6, "request network map every n hours")
	flag.Parse()

	if *broker == "" || *user == "" || *pass == "" || *basetopic == "" {
		flag.Usage()
		return
	}
	
	go run_mqtt(*broker, *user, *pass, *basetopic, *mapNHrs)
	
	marshaller = declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	serve(*addr)
}
