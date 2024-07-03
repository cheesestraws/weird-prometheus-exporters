package main

import (
	"os"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"net/http"
	"log"
	"time"
	
	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

var metrics Metrics

func mangleParams() ([]int, error) {
	sp := flag.String("stations", "43160,1754,43159,43165,43159", "comma-separated list of station IDs; see https://environment.data.gov.uk/flood-monitoring/id/stations/")
	flag.Parse()

	if *sp == "" {
		return nil, errors.New("please give one or more station IDs")
	}

	strs := strings.Split(*sp, ",")
	var ints []int

	for _, v := range strs {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("given station %q is not a number", v)
		}
		ints = append(ints, i)
	}

	return ints, nil
}

func runloop(cli *http.Client, sleepTime time.Duration, stationIDs []int) {
	for {
		stations, err := fn.Errmap(stationIDs, func(stationID int) (*station, error) {
			return fetch(cli, stationID)
		})
		if err != nil {
			// Don't propagate the error upwards, print it and carry on
			log.Printf("err: %v", err)
			break
		}
	
		levels := fn.Map(stations, riverLevelFromStation)
		metrics.Import(levels)
		fmt.Printf("%s", metrics.ToPrometheus())
	
		time.Sleep(sleepTime)
	}
}


func main() {
	stations, err := mangleParams()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	client := makeHTTPClient()
	runloop(client, 10 * time.Minute, stations)
}
