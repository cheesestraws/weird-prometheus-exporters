package main

import (
	"time"
)

type Summaries struct {
	Window time.Duration

	NumTrains int
	NumLateTrains int
	NumCancelledTrains int
	
	AvgLateTime time.Duration
	WorstLateTime time.Duration
	
	CancelReasons map[string]int
}

func (ss WrappedServices) Summarise(window time.Duration) *Summaries {
	sum := &Summaries{
		Window: window,
		CancelReasons: make(map[string]int),
	}
	
	var latenessAccumulator time.Duration
	
	for _, s := range ss {
		if !s.Valid {
			continue
		}
		
		sum.NumTrains++
		
		if s.Lateness > 5 * time.Minute {
			sum.NumLateTrains++
		}
		
		if s.S.LocationDetail.DisplayAs == "CANCELLED_CALL" {
			sum.NumCancelledTrains++
		}
		
		latenessAccumulator += s.Lateness
		
		if s.Lateness > sum.WorstLateTime {
			sum.WorstLateTime = s.Lateness
		}
		
		if s.S.LocationDetail.CancelReasonShortText != "" {
			sum.CancelReasons[s.S.LocationDetail.CancelReasonShortText]++
		} else if s.S.LocationDetail.CancelReasonLongText != "" {
			sum.CancelReasons[s.S.LocationDetail.CancelReasonLongText]++
		} else if s.S.LocationDetail.CancelReasonCode != "" {
			sum.CancelReasons[s.S.LocationDetail.CancelReasonCode]++
		}
		
	}
	
	return sum
}