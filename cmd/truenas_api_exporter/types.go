package main

import (
	"regexp"
	"strconv"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
	"github.com/cheesestraws/weird-prometheus-exporters/lib/truenas"
)

type PoolIdentifier struct {
	ID   int
	Name string
	Path string
}

type CloudSyncIdentifier struct {
	ID          int
	Description string
	Path        string
}

type Summary struct {
	ActiveAlertCount map[string]int

	PoolHealthy        map[PoolIdentifier]bool
	PoolStatus         map[PoolIdentifier]int
	PoolReadErrors     map[PoolIdentifier]int
	PoolWriteErrors    map[PoolIdentifier]int
	PoolChecksumErrors map[PoolIdentifier]int
	PoolSize           map[PoolIdentifier]int
	PoolAllocated      map[PoolIdentifier]int

	CloudSyncState          map[CloudSyncIdentifier]int
	CloudSyncAllegedPercent map[CloudSyncIdentifier]int
	CloudSyncDoneBytes      map[CloudSyncIdentifier]int
}

var poolStatusMap = map[string]int{
	"unknown":  0,
	"ONLINE":   1,
	"DEGRADED": 2,
	"FAULTED":  3,
	"OFFLINE":  4,
	"UNAVAIL":  5,
	"REMOVED":  6,
}

var syncStatusMap = map[string]int{
	"unknown": 0,
	"FAILED":  1,
	"ABORTED": 2,
	"PENDING": 3,
	"RUNNING": 4,
	"SUCCESS": 5,
}

var quantityRE = regexp.MustCompile("^([0-9.]+)(K|M|G|T)")

func stringToBytes(s string) (int, bool) {
	multipliers := map[string]float64{
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
		"T": 1024 * 1024 * 1024 * 1024,
	}

	matches := quantityRE.FindStringSubmatch(s)
	if matches == nil {
		return 0, false
	}

	if len(matches) != 3 {
		return 0, false
	}

	f, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, false
	}

	total := f * multipliers[matches[2]]

	if total > 0 {
		return int(total), true
	} else {
		return 0, false
	}
}

func summarise(alerts []truenas.Alert, pools []truenas.Pool, syncs []truenas.CloudSync) Summary {
	sum := Summary{
		ActiveAlertCount: make(map[string]int),

		PoolHealthy:        make(map[PoolIdentifier]bool),
		PoolStatus:         make(map[PoolIdentifier]int),
		PoolReadErrors:     make(map[PoolIdentifier]int),
		PoolWriteErrors:    make(map[PoolIdentifier]int),
		PoolChecksumErrors: make(map[PoolIdentifier]int),
		PoolSize:           make(map[PoolIdentifier]int),
		PoolAllocated:      make(map[PoolIdentifier]int),

		CloudSyncState:          make(map[CloudSyncIdentifier]int),
		CloudSyncAllegedPercent: make(map[CloudSyncIdentifier]int),
		CloudSyncDoneBytes:      make(map[CloudSyncIdentifier]int),
	}

	// Alerts first.  Zeroes to begin with
	for _, severity := range truenas.AlertValidLevels {
		sum.ActiveAlertCount[severity] = 0
	}

	activeAlerts := fn.Filter(alerts, func(a truenas.Alert) bool {
		return !a.Dismissed
	})

	for _, a := range activeAlerts {
		sum.ActiveAlertCount[a.Level]++
	}

	// Pools
	for _, p := range pools {
		id := PoolIdentifier{
			ID:   p.ID,
			Name: p.Name,
			Path: p.Path,
		}

		sum.PoolHealthy[id] = p.Healthy
		sum.PoolStatus[id] = poolStatusMap[p.Status]

		// Sum stats from top-level data devices
		readErrors := 0
		writeErrors := 0
		checksumErrors := 0
		size := 0
		allocated := 0
		for _, d := range p.Topology.Data {
			readErrors += d.Stats.ReadErrors
			writeErrors += d.Stats.WriteErrors
			checksumErrors += d.Stats.ChecksumErrors
			size += d.Stats.Size
			allocated += d.Stats.Allocated
		}

		sum.PoolReadErrors[id] = readErrors
		sum.PoolWriteErrors[id] = writeErrors
		sum.PoolChecksumErrors[id] = checksumErrors
		sum.PoolSize[id] = size
		sum.PoolAllocated[id] = allocated
	}

	// Cloud syncs
	for _, s := range syncs {
		id := CloudSyncIdentifier{
			ID:          s.ID,
			Description: s.Description,
			Path:        s.Path,
		}

		sum.CloudSyncState[id] = syncStatusMap[s.Job.State]
		sum.CloudSyncAllegedPercent[id] = s.Job.Progress.Percent

		b, ok := stringToBytes(s.Job.Progress.Description)
		if ok {
			sum.CloudSyncDoneBytes[id] = b
		}
	}

	return sum
}
