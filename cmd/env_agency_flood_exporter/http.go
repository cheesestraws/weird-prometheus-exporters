package main

import (
	"net/http"
	"encoding/json"
	"io"
	"fmt"
	"time"
	"log"
)

func makeHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func fetch(cli *http.Client, stationID int) (*station, error) {
	url := fmt.Sprintf("https://environment.data.gov.uk/flood-monitoring/id/stations/%d", stationID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sta station
	err = json.Unmarshal(body, &sta)
	if err != nil {
		return nil, err
	}
	
	return &sta, nil
}

func serve(addr string, stationIDs []int, m *Metrics) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(m.ToPrometheus())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "polling stations %v", stationIDs)
	})
	
	log.Printf("listening on %s", addr)
	
	log.Fatal(http.ListenAndServe(addr, nil))
}