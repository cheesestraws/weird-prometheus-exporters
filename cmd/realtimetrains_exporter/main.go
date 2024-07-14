package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	rtt "github.com/cheesestraws/gortt"
)

var state struct {
	station *string
	dump    *bool

	l     sync.RWMutex
	stuff []byte
}

func doFetchState() {
	var summarieses []*Summaries

	cli, err := rtt.NewClient("", "")
	if err != nil {
		log.Printf("rtt.NewClient: %v", err)
		return
	}

	// Make a consistent snapshots of the time windows
	windows := config.TimeWindows.Snapshot()

	furthestBack := time.Unix(math.MaxInt64/2, 0)
	furthestFwd := time.Time{}

	for _, w := range windows {
		if w.From.Before(furthestBack) {
			furthestBack = w.From
		}

		if w.To.After(furthestFwd) {
			furthestFwd = w.To
		}
	}

	fetches := MakeFetches(*state.station, furthestBack, furthestFwd)
	allServices, err := fetches.Do(context.Background(), cli)
	if err != nil {
		log.Printf("fetches.Do: %v", err)
	}

	for _, window := range config.TimeWindows {
		services := allServices.ByTimeWindow(window.From(), window.To())

		// "All trains" summary
		sums := services.Summarise(window.Name)
		summarieses = append(summarieses, sums)

		// Per destination
		for _, l := range services.Destinations() {
			loc := l // copy here because we need to take ref later
			filtered := services.ByDestinationTIPLOC(l.TIPLOC)
			sums := filtered.Summarise(window.Name)
			sums.Destination = &loc

			summarieses = append(summarieses, sums)
		}

		// Per origin
		for _, l := range services.Origins() {
			loc := l // copy here because we need to take ref later
			filtered := services.ByOriginTIPLOC(l.TIPLOC)
			sums := filtered.Summarise(window.Name)
			sums.Origin = &loc

			summarieses = append(summarieses, sums)
		}

	}

	// write em out and sell em cheap
	bs := &bytes.Buffer{}
	for _, v := range summarieses {
		bs.Write(v.Prometheise())
		bs.Write([]byte("\n\n"))
	}

	if *state.dump {
		log.Printf("%s\n", bs)
	}

	state.l.Lock()
	state.stuff = bs.Bytes()
	state.l.Unlock()
}

func fetch_and_update_forever() {
	for {
		doFetchState()
		time.Sleep(2 * time.Minute)
	}
}

func serve(addr string) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		var bs []byte
		state.l.RLock()
		bs = state.stuff
		state.l.RUnlock()

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
	addr := flag.String("addr", ":9403", "address to listen on")
	state.station = flag.String("station", "SAL", "station to watch trains at")
	state.dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	flag.Parse()

	go fetch_and_update_forever()
	serve(*addr)
}
