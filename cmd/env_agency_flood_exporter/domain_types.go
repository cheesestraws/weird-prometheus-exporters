package main

// domain_types.go contains types that contain the info we actually care
// about

import (
	"strings"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

type RiverLevel struct {
	Valid bool
	When  time.Time

	StationID    string
	StationLabel string
	RiverName    string
	Level        float64

	typicalHigh *float64
	typicalLow  *float64
}

func (r RiverLevel) TypicalHigh(fudges Fudges) *float64 {
	_, hi, ok := fudges.Get(r.StationID)
	if ok {
		return &hi
	}
	
	return r.typicalHigh
}

func (r RiverLevel) TypicalLow(fudges Fudges) *float64 {
	lo, _, ok := fudges.Get(r.StationID)
	if ok {
		return &lo
	}
	
	return r.typicalLow
}


var trimLabelPrefix string

func riverLevelFromStation(s *station) RiverLevel {
	var r RiverLevel

	if s == nil {
		return r
	}

	r.StationID = s.Items.Notation

	l := strings.ToLower(s.Items.Label)
	trimpfx := strings.ToLower(trimLabelPrefix)
	l = strings.TrimPrefix(l, trimpfx)

	r.StationLabel = strings.Title(l)
	r.RiverName = s.Items.RiverName

	// Now, do we have a level measurement?
	ms := fn.Filter(s.Items.Measures, func(m measure) bool {
		return m.Qualifier == "Stage" && m.Parameter == "level" &&
			m.LatestReading.Valid
	})

	if len(ms) == 0 {
		return r
	}

	if ms[0].LatestReading.Timestamp == nil {
		return r
	}

	r.When = time.Time(*ms[0].LatestReading.Timestamp)
	r.Level = ms[0].LatestReading.Value

	// At this point we're valid
	r.Valid = true

	// Do we have our bonus information?
	if s.Items.StageScale != nil {
		r.typicalHigh = s.Items.StageScale.TypicalRangeHigh
		r.typicalLow = s.Items.StageScale.TypicalRangeLow
	}

	return r
}
