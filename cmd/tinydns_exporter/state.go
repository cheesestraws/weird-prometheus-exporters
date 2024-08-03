package main

import (
	"bytes"
	"context"
	"log"
	"sync"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

var state struct {
	l sync.RWMutex

	data fn.Maybe[DatabaseSummary]
	cdb  fn.Maybe[CDBStats]
	svc  fn.Maybe[ServiceStatus]

	errs errors
}

type errors struct {
	dataErrors int `prometheus:"data_error_count"`
	cdbErrors  int `prometheus:"cdb_error_count"`
	svcErrors  int `prometheus:"svstat_error_count"`
}

func start() {
	if *datafile != "" {
		go dataRunloop()
	}

	if *datacdb != "" {
		go cdbRunloop()
	}

	if *servicedir != "" {
		go svcRunloop()
	}
}

func dataRunloop() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		ds, err := checkData(ctx, *datafile, *suffix)

		state.l.Lock()
		if err != nil {
			log.Printf("data runloop: %v", err)
			state.data = fn.Absent[DatabaseSummary]()
			state.errs.dataErrors++
		} else {
			state.data = fn.Present(ds)
		}
		state.l.Unlock()

		time.Sleep(5 * time.Minute)
	}
}

func cdbRunloop() {
	for {
		st, err := checkCDB(*datacdb)
		state.l.Lock()
		if err != nil {
			log.Printf("cdb runloop: %v", err)
			state.cdb = fn.Absent[CDBStats]()
			state.errs.cdbErrors++
		} else {
			state.cdb = fn.Present(st)
		}
		state.l.Unlock()

		time.Sleep(1 * time.Minute)
	}
}

func svcRunloop() {
	for {
		sv, err := Svstat(*servicedir)
		state.l.Lock()
		if err != nil {
			log.Printf("svstat runloop: %v", err)
			state.svc = fn.Absent[ServiceStatus]()
			state.errs.svcErrors++
		} else {
			state.svc = fn.Present(sv)
		}
		state.l.Unlock()

		time.Sleep(1 * time.Minute)
	}
}

func getBody() []byte {
	var bs bytes.Buffer
	m := declprom.Marshaller{
		MetricNamePrefix: "tinydns_",
	}

	state.l.RLock()
	defer state.l.RUnlock()

	state.data.Range(func(ds DatabaseSummary) {
		bs.Write(m.Marshal(ds, nil))
	})

	state.cdb.Range(func(cdb CDBStats) {
		bs.Write(m.Marshal(cdb, nil))
	})

	state.svc.Range(func(svc ServiceStatus) {
		bs.Write(m.Marshal(svc, nil))
	})

	bs.Write(m.Marshal(state.errs, nil))

	return bs.Bytes()
}
