package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
	"github.com/cheesestraws/weird-prometheus-exporters/lib/truenas"
)

var baseURL *string
var user *string
var pass *string
var prefix *string
var addr *string
var dump *bool

func fetch() ([]byte, error) {
	cli := truenas.NewClient(*baseURL, *user, *pass, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alerts, err := cli.AlertList(ctx)
	if err != nil {
		return nil, err
	}

	pools, err := cli.Pools(ctx)
	if err != nil {
		return nil, err
	}

	syncs, err := cli.CloudSyncs(ctx)
	if err != nil {
		return nil, err
	}

	sums := summarise(alerts, pools, syncs)

	m := declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	bs := m.Marshal(sums, map[string]string{
		"base_url": *baseURL,
	})

	if *dump {
		log.Printf("%+s", bs)
	}

	return bs, nil
}

func serve(addr string) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		bs, _ := fetch()
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
	baseURL = flag.String("baseurl", os.Getenv("NASURL"), "base URL of Lyrion instance")
	user = flag.String("user", os.Getenv("NASUSER"), "user name")
	pass = flag.String("pass", os.Getenv("NASPASS"), "password")
	prefix = flag.String("prefix", "truenas_", "prefix for metric names")
	addr = flag.String("addr", ":9408", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	flag.Parse()

	if *dump {
		fetch()
	}

	serve(*addr)
}
