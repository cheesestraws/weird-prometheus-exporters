package main

import (
	"sync"
	"log"
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

func main() {
	log.Printf("initial state fetch...")
	doFetchState()
	
	log.Printf("%s", state.stuff)
}