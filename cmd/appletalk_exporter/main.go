package main

import (
	"sync"
	"log"
	"time"
	"net/http"
)

var state struct {
	l sync.RWMutex
	stuff []byte
}

func doFetchState() {
	ns, err := QueryNetworkState()
	if err != nil {
		log.Printf("%v", err)
	}
	
	bs := ns.ToPrometheus()
	state.l.Lock()
	state.stuff=bs
	state.l.Unlock()
}

func fetch_and_update_forever() {
	for {
		doFetchState()
		time.Sleep(30 * time.Second)
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
	addr := flag.String("addr", ":9402", "address to listen on")
	flag.Parse()	
		
	go fetch_and_update_forever()
	serve(*addr)
	
	
}