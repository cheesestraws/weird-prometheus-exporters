package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"time"
)

var alwaysLabelsRegardlessOfType map[string]struct{} = map[string]struct{}{
	"id":      struct{}{},
	"channel": struct{}{},
}

func streamEvents(baseURL string) {
	url, err := url.JoinPath(baseURL, "/stream")
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(1 * time.Second)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("stream: http.Get: %v", err)
		}

		err = handleConn(resp.Body)
		if err != nil {
			log.Printf("stream: handleConn: %v", err)
		}
	}
}

func handleConn(c io.ReadCloser) error {
	log.Printf("HTTP stream up")
	defer log.Printf("HTTP stream closed")
	defer c.Close()

	metrics.Lock()
	metrics.StreamConnectionUp = 1
	metrics.Unlock()

	defer func() {
		metrics.Lock()
		metrics.StreamConnectionUp = 0
		metrics.Unlock()
	}()

	dec := json.NewDecoder(c)
	for {
		mm := make(map[string]JSONNumberOrString)
		err := dec.Decode(&mm)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		err = handleLine(mm)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleLine(mm map[string]JSONNumberOrString) error {
	// Check that our attributes only contain numbers and strings
	for _, val := range mm {
		// Occasionally rtl_433 sends us an unsolicited stats
		// message which contains an array.  It's not clear what
		// we should do with arrays and other crap so
		// let's just filter them out and hope for the best
		if val.IsOther {
			log.Printf("got unexpected field type; errant stats request?")
			return nil
		}
	}

	// messages from entities have a model and time; anything else,
	// ignore
	_, hasTime := mm["time"]
	_, hasModel := mm["model"]

	if !hasTime || !hasModel {
		log.Printf("unexpecyed: %+v", mm)
		return nil
	}

	fields := maps.Clone(mm)

	delete(fields, "time")

	// make label set and metric set; labels are strings,
	// metrics are numbers
	labelSet := make(map[string]string)
	metricSet := make(map[string]float64)
	for k, v := range fields {
		if _, ok := alwaysLabelsRegardlessOfType[k]; ok {
			labelSet[k] = v.LabelString()
		} else if v.IsNumber {
			metricSet[k] = v.Number
		} else if v.IsString {
			labelSet[k] = v.String
		}
	}

	registerMetric(labelSet, metricSet)

	if *dump {
		fmt.Printf("%s\n\n\n", getBody())
	}

	return nil
}

func registerMetric(labels map[string]string, values map[string]float64) {
	metrics.Lock()
	defer metrics.Unlock()

	metrics.DynamicMetrics.Observe(labels, values)
}
