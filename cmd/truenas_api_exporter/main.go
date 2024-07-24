package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/truenas"
)

var baseURL *string
var user *string
var pass *string
var prefix *string
var addr *string
var dump *bool

var state struct {
	l    sync.RWMutex
	sums Summary
}

func fetch() error {
	cli := truenas.NewClient(*baseURL, *user, *pass, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alerts, err := cli.AlertList(ctx)
	if err != nil {
		return err
	}

	pools, err := cli.Pools(ctx)
	if err != nil {
		return err
	}

	syncs, err := cli.CloudSyncs(ctx)
	if err != nil {
		return err
	}

	state.l.Lock()
	defer state.l.Unlock()
	state.sums = summarise(alerts, pools, syncs)

	return nil
}

func main() {
	baseURL = flag.String("baseurl", os.Getenv("NASURL"), "base URL of Lyrion instance")
	user = flag.String("user", os.Getenv("NASUSER"), "user name")
	pass = flag.String("pass", os.Getenv("NASPASS"), "password")
	prefix = flag.String("prefix", "truenas_", "prefix for metric names")
	addr = flag.String("addr", ":9408", "address to listen on")
	dump = flag.Bool("d", false, "dump metrics to stdout as well as http")
	flag.Parse()

	fetch()

	log.Printf("%+v", state.sums)
}
