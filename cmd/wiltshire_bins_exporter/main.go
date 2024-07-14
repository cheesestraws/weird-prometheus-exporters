package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
	"github.com/cheesestraws/weird-prometheus-exporters/lib/wiltshirebins"
)

var state struct {
	dump     *bool
	postcode *string
	uprn     *string

	l sync.RWMutex

	stuff []byte
}

// truncdate, truncates to the nearest local day.  it's a pun, geddit
func truncdate(t time.Time) time.Time {
	yyyy, mm, dd := t.Date()

	return time.Date(yyyy, mm, dd, 0, 0, 0, 0, time.Local)
}

func fetch_and_update_forever() {
	m := declprom.Marshaller{
		MetricNamePrefix: "bincollection_",
		MetricNameSuffix: "_tomorrow",
		BaseLabels: map[string]string{
			"postcode": *state.postcode,
			"uprn":     *state.uprn,
		},
	}

	for {
		today := truncdate(time.Now())
		tomorrow := today.Add(26 * time.Hour)
		
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		collection, err := wiltshirebins.DefaultClient.GetForDate(ctx, tomorrow, *state.postcode, *state.uprn)
		if err != nil {
			log.Printf("fetch error: %v", err)
		}

		bs := m.Marshal(collection, nil)

		// stash the bytes in the state
		state.l.Lock()
		state.stuff = bs
		state.l.Unlock()

		if *state.dump {
			log.Printf("%s", bs)
		}

		time.Sleep(3 * time.Hour)
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
	// cli arguments
	addr := flag.String("addr", ":9404", "address to listen on")
	state.dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	state.postcode = flag.String("postcode", "", "postcode of your house")
	state.uprn = flag.String("uprn", "", "the unique property reference number of your house; you can find this at https://www.findmyaddress.co.uk")
	flag.Parse()

	if *state.postcode == "" {
		log.Printf("postcode must be provided")
		return
	}

	if *state.uprn == "" {
		log.Printf("uprn must be provided")
		return
	}

	go fetch_and_update_forever()

	serve(*addr)
}
