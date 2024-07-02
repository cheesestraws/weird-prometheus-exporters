package main

// domain_types.go contains types that contain the info we actually care 
// about

import (
	"time"
	
	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

type RiverLevel struct {
	Valid bool
	When time.Time

	StationID string
	StationLabel string
	RiverName string
	Level float64
	
	TypicalHigh *float64
	TypicalLow *float64
}

func riverLevelFromStation(s *station) RiverLevel {
	var r RiverLevel
	
	if s == nil {
		return r
	}
	
	r.StationID = s.Items.Notation
	r.StationLabel = s.Items.Label
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
		r.TypicalHigh = s.Items.StageScale.TypicalRangeHigh
		r.TypicalLow = s.Items.StageScale.TypicalRangeLow
	}
	
	return r
}