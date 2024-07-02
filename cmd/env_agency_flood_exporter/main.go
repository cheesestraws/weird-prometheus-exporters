package main

import (
	"os"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	
	_ "github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

func mangleParams() ([]int, error) {
	sp := flag.String("stations", "", "comma-separated list of station IDs; see https://environment.data.gov.uk/flood-monitoring/id/stations/")
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


func main() {
	/*stations, err := mangleParams()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}*/

	client := makeHTTPClient()
	err := fetch(client, 43159)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
