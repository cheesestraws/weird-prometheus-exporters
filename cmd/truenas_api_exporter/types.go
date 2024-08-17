package main

import (
	"regexp"
	"strconv"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
	"github.com/cheesestraws/weird-prometheus-exporters/lib/truenas"
)

type PoolIdentifier struct {
	ID   int    `prometheus_label:"pool_id"`
	Name string `prometheus_label:"pool_name"`
	Path string `prometheus_label:"pool_path"`
}

type CloudSyncIdentifier struct {
	ID          int    `prometheus_label:"job_id"`
	Description string `prometheus_label:"job_description"`
	Path        string `prometheus_label:"job_path"`
}

type Summary struct {
	ActiveAlertCount map[string]int `prometheus_map:"active_alerts" prometheus_map_key:"severity"`

	PoolHealthy        map[PoolIdentifier]int `prometheus_map:"pool_healthy"`
	PoolStatus         map[PoolIdentifier]int `prometheus_map:"pool_status" prometheus_help:"0 unknown 1 ONLINE 2 DEGRADED 3 FAULTED 4 OFFLINE 5 UNAVAIL 6 REMOVED"`
	PoolReadErrors     map[PoolIdentifier]int `prometheus_map:"pool_read_errors"`
	PoolWriteErrors    map[PoolIdentifier]int `prometheus_map:"pool_write_errors"`
	PoolChecksumErrors map[PoolIdentifier]int `prometheus_map:"pool_checksum_errors"`
	PoolSize           map[PoolIdentifier]int `prometheus_map:"pool_size"`
	PoolAllocated      map[PoolIdentifier]int `prometheus_map:"pool_allocated"`
	PoolPctAllocated   map[PoolIdentifier]int `prometheus_map:"pool_allocated_pct"`

	CloudSyncState          map[CloudSyncIdentifier]int `prometheus_map:"cloudsync_state" prometheus_help:"0 unknown 1 FAILED 2 ABORTED 3 PENDING 4 RUNNING 5 SUCCESS"`
	CloudSyncEnabled        map[CloudSyncIdentifier]int `prometheus_map:"cloudsync_enabled"`
	CloudSyncAllegedPercent map[CloudSyncIdentifier]int `prometheus_map:"cloudsync_alleged_progress"`
	CloudSyncDoneBytes      map[CloudSyncIdentifier]int `prometheus_map:"cloudsync_done_bytes"`
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

		PoolHealthy:        make(map[PoolIdentifier]int),
		PoolStatus:         make(map[PoolIdentifier]int),
		PoolReadErrors:     make(map[PoolIdentifier]int),
		PoolWriteErrors:    make(map[PoolIdentifier]int),
		PoolChecksumErrors: make(map[PoolIdentifier]int),
		PoolSize:           make(map[PoolIdentifier]int),
		PoolAllocated:      make(map[PoolIdentifier]int),
		PoolPctAllocated:   make(map[PoolIdentifier]int),

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

		if p.Healthy {
			sum.PoolHealthy[id] = 1
		} else {
			sum.PoolHealthy[id] = 0
		}
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

		sum.PoolPctAllocated[id] = int((float64(allocated) / float64(size)) * 100)
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
