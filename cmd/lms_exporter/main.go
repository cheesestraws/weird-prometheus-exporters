package main

import (
	"context"
	"flag"
	"log"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

var baseURL *string
var prefix *string
var addr *string
var dump *bool

func getBody() ([]byte, error) {
	ctx := context.Background()
	s, err := fetchStatus(ctx, *baseURL)
	if err != nil {
		return nil, err
	}

	m := declprom.Marshaller{
		MetricNamePrefix: *prefix,
	}

	return m.Marshal(s.Summarise(), nil), nil
}

func main() {
	baseURL = flag.String("baseurl", "http://music.lan:9000", "base URL of Lyrion instance")
	prefix = flag.String("prefix", "lms_", "prefix for metric names")
	addr = flag.String("addr", ":9406", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	flag.Parse()

	b, err := getBody()
	if err != nil {
		log.Printf("err: %v", err)
	}

	log.Printf("%s", b)

}
