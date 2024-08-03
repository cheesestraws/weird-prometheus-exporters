package main

import (
	"context"
	"log"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

func main() {
	ctx := context.Background()
	dataSummary, err := checkData(ctx, "/Users/cheesey/data", ".ns.ecliptiq.co.uk")
	if err != nil {
		log.Printf("err: %v", err)
	}

	m := declprom.Marshaller{
		MetricNamePrefix: "tinydns_",
	}
	bs := m.Marshal(dataSummary, nil)
	log.Printf("%s\n", bs)
}
