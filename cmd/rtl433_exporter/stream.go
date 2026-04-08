package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/rflookup"
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

var autoLevelRegexp = regexp.MustCompile(`Estimated noise level is (-?[0-9\.]+) dB, adjusting minimum detection level to (-?[0-9\.]+) dB`)
func handleAutoLevelLine(mm map[string]JSONNumberOrString) bool {
	src, hasSrc := mm["src"]
	if !hasSrc {
		return false
	}
	
	if !src.IsString {
		return false
	}
	
	if src.String != "Auto Level" {
		return false
	}
	
	msg, hasMsg := mm["msg"]
	if !hasMsg {
		return false
	}
	
	if !msg.IsString {
		return false
	}
	
	ss := autoLevelRegexp.FindStringSubmatch(mm["msg"].String)
	if len(ss) < 3 {
		return false
	}
	
	noise, err := strconv.ParseFloat(ss[1], 64)
	if err != nil {
		log.Printf(`error parsing line "%+v": %s should have been a number`, mm, ss[1])
		return false
	}
	
	sens, err := strconv.ParseFloat(ss[2], 64)
	if err != nil {
		log.Printf(`error parsing line "%+v": %s should have been a number`, mm, ss[2])
		return false
	}

	metrics.Lock()
	defer metrics.Unlock()
	
	metrics.EstimatedNoiseLevel = noise
	metrics.MinimumDetectionLevel = sens
	
	return true
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

	// Auto level log messages are their own beasts and are a bit icky
	if handleAutoLevelLine(mm) {
		return nil
	}

	// messages from entities have a model and time; anything else,
	// ignore
	_, hasTime := mm["time"]
	_, hasModel := mm["model"]

	if !hasTime || !hasModel {
		log.Printf("unexpected: %+v", mm)
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

	if *useRFLookup {
		handleRFLookup(labelSet)
	}

	registerMetric(labelSet, metricSet)

	if *dump {
		fmt.Printf("%s\n\n\n", getBody())
	}

	return nil
}

func handleRFLookup(labelset map[string]string) {
	// Do we have metadata yet?
	centreFrequency := metrics.GetMetadata().CentreFrequency
	if centreFrequency == 0 {
		return
	}

	rflookupName, err := rflookup.MkLookupHostname(labelset, centreFrequency)
	if err != nil {
		metrics.Lock()
		metrics.RFLookupErrors++
		metrics.Unlock()
		log.Printf("rflookup MkLookupHostname: $v", err)
		return
	}

	labelset["rflookup"] = rflookupName

	metrics.Lock()
	metrics.RFLookupTotal++
	metrics.Unlock()

	cname, err := rflookup.Lookup(labelset, centreFrequency)
	if err == nil {
		metrics.Lock()
		metrics.RFLookupHits++
		metrics.Unlock()

		labelset["rflookup_friendly"] = cname
	} else if err != nil && strings.Contains(err.Error(), "no such host") {
		metrics.Lock()
		metrics.RFLookupAnonymous++
		metrics.Unlock()
	} else {
		metrics.Lock()
		metrics.RFLookupErrors++
		metrics.Unlock()
		log.Printf("rflookup Lookup: %v", err)
	}
}

func registerMetric(labels map[string]string, values map[string]float64) {
	metrics.Lock()
	defer metrics.Unlock()

	metrics.DynamicMetrics.Observe(labels, values)
}
