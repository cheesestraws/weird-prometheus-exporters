package main

import (
	"bytes"
	"fmt"
	"sync"
)

type Metrics struct {
	l sync.RWMutex

	m map[string]RiverLevel
}

func (m *Metrics) set(r RiverLevel) {
	// keep the zero value usable
	if m.m == nil {
		m.m = make(map[string]RiverLevel)
	}
	m.m[r.StationID] = r
}

func (m *Metrics) Import(rs []RiverLevel) {
	m.l.Lock()
	defer m.l.Unlock()

	for _, r := range rs {
		if r.Valid {
			m.set(r)
		}
	}
}

func (m *Metrics) ToPrometheus() []byte {
	m.l.RLock()
	defer m.l.RUnlock()

	var b bytes.Buffer
	const pfx = "environment_data_gov_uk_flood_monitoring"

	writeDimensions := func(r RiverLevel) {
		fmt.Fprintf(&b, "{station_id=%q,station_label=%q", r.StationID, r.StationLabel)

		if r.RiverName != "" {
			fmt.Fprintf(&b, ",river_name=%q", r.RiverName)
		}

		fmt.Fprintf(&b, "}")
	}

	var fst bool

	// Do the level metrics
	fst = true
	for _, r := range m.m {
		if fst {
			fmt.Fprintf(&b, "# HELP %s_river_level The level of the river.\n", pfx)
			fmt.Fprintf(&b, "# TYPE %s_river_level gauge\n", pfx)
			fst = false
		}
		fmt.Fprintf(&b, "%s_river_level", pfx)
		writeDimensions(r)
		fmt.Fprintf(&b, " %v %d\n", r.Level, r.When.UnixMilli())
	}

	fmt.Fprintf(&b, "\n\n")

	// Typical highs
	fst = true
	for _, r := range m.m {
		if r.TypicalHigh(fudges) != nil {
			if fst {
				fmt.Fprintf(&b, "# HELP %s_river_typical_high The highest level expected under normal circumstances.\n", pfx)
				fmt.Fprintf(&b, "# TYPE %s_river_typical_high gauge\n", pfx)
				fst = false
			}

			fmt.Fprintf(&b, "%s_river_typical_high", pfx)
			writeDimensions(r)
			fmt.Fprintf(&b, " %v %d\n", *r.TypicalHigh(fudges), r.When.UnixMilli())
		}
	}

	fmt.Fprintf(&b, "\n\n")

	// Typical lows
	fst = true
	for _, r := range m.m {
		if r.TypicalLow(fudges) != nil {
			if fst {
				fmt.Fprintf(&b, "# HELP %s_river_typical_low The lowest level expected under normal circumstances.\n", pfx)
				fmt.Fprintf(&b, "# TYPE %s_river_typical_low gauge\n", pfx)
				fst = false
			}

			fmt.Fprintf(&b, "%s_river_typical_low", pfx)
			writeDimensions(r)
			fmt.Fprintf(&b, " %v %d\n", *r.TypicalLow(fudges), r.When.UnixMilli())
		}
	}

	return b.Bytes()
}
