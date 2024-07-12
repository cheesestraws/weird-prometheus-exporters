package main

import (
	"context"
	"log"
	"time"

	rtt "github.com/cheesestraws/gortt"
)

func main() {
	cli, err := rtt.NewClient("", "")
	if err != nil {
		log.Printf("rtt.NewClient: %v", err)
		return
	}

	from := time.Now().Add(-1 * time.Hour)
	to := time.Now().Add(30 * time.Minute)

	fetches := MakeFetches("SAL", from, to)
	services, err := fetches.Do(context.Background(), cli)
	if err != nil {
		log.Printf("fetches.Do: %v", err)
	}

	services = services.ByTimeWindow(from, to)
	sums := services.Summarise(90*time.Minute, "SAL")

	log.Printf("")
	log.Printf("ok")
	log.Printf("")

	log.Printf("summary: %+v", sums)
		
	log.Printf("")
	log.Printf("prometheus:")
	log.Printf("\n%s", sums.Prometheise())

}
