package main

import (
	"context"
	"log"
	"time"
	"sync"
	"fmt"
	"flag"
	"net/http"

	rtt "github.com/cheesestraws/gortt"
)

var state struct {
	station *string
	dump *bool

	l     sync.RWMutex	
	stuff []byte
}

func doFetchState() {
	cli, err := rtt.NewClient("", "")
	if err != nil {
		log.Printf("rtt.NewClient: %v", err)
		return
	}

	from := time.Now().Add(-1 * time.Hour)
	to := time.Now().Add(30 * time.Minute)

	fetches := MakeFetches(*state.station, from, to)
	services, err := fetches.Do(context.Background(), cli)
	if err != nil {
		log.Printf("fetches.Do: %v", err)
	}

	services = services.ByTimeWindow(from, to)
	sums := services.Summarise(90*time.Minute)
	
	bs := sums.Prometheise()
	
	if *state.dump {
		log.Printf("%s\n", bs)
	}
	
	log.Printf("Destinations: %+v", services.Destinations())
	log.Printf("Origins: %+v", services.Origins())

	
	state.l.Lock()
	state.stuff = bs
	state.l.Unlock()
}

func fetch_and_update_forever() {
	for {
		doFetchState()
		time.Sleep(2 *  time.Minute)
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
