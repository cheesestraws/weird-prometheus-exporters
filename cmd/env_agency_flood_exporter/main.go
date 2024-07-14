package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

var metrics Metrics
var fudges Fudges

func mangleParams() ([]int, string, error) {
	var err error

	sp := flag.String("stations", "43160,1754,43159,43165", "comma-separated list of station IDs; see https://environment.data.gov.uk/flood-monitoring/id/stations/")
	addr := flag.String("addr", ":9401", "address to listen on")
	flag.StringVar(&trimLabelPrefix, "strip-prefix", "salisbury ", "prefix to strip from station labels")
	fudgeStrings := flag.String("fudges", "1754:0.35:0.84", "manually fudge typical high and low for stations that don't have them or where they are wrong.  A comma-separated list of station:hi:lo triples.")
	flag.Parse()

	if *sp == "" {
		return nil, "", errors.New("please give one or more station IDs")
	}

	strs := strings.Split(*sp, ",")
	var ints []int

	for _, v := range strs {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, "", fmt.Errorf("given station %q is not a number", v)
		}
		ints = append(ints, i)
	}

	fudges, err = ParseFudges(*fudgeStrings)
	if err != nil {
		return ints, *addr, err
	}

	return ints, *addr, nil
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

		time.Sleep(sleepTime)
	}
}

func main() {
	stations, addr, err := mangleParams()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	client := makeHTTPClient()
	go runloop(client, 10*time.Minute, stations)

	serve(addr, stations, &metrics)
}
