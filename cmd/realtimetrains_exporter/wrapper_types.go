package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"

	rtt "github.com/cheesestraws/gortt"
)

type WrappedService struct {
	Valid bool

	RequestDate time.Time

	GBTTDeparture     time.Time
	RealtimeDeparture time.Time
	Lateness          time.Duration

	S rtt.RTTLocationContainer
}

type WrappedServices []WrappedService

func trainTimeToGoTime(date time.Time, t string) (time.Time, error) {
	if len(t) != 4 {
		return time.Time{}, fmt.Errorf("%q is a bad time", t)
	}

	h := t[0:2]
	m := t[2:4]
	dur := fmt.Sprintf("%sh%sm", h, m)
	d, err := time.ParseDuration(dur)
	if err != nil {
		return time.Time{}, err
	}

	ret := date.Add(d)
	return ret, nil
}

// Converts a Location lineup from RTT into a slice of WrappedServices
// ready for filtering, dicing, slicing and various mangling.
func LocationLineupToServices(ll rtt.RTTLocationLineup, date time.Time) WrappedServices {
	return WrappedServices(fn.Map(ll.Services, func(s rtt.RTTLocationContainer) WrappedService {
		valid := true
		gbttDeparture, err := trainTimeToGoTime(date, s.LocationDetail.GBTTBookedDeparture)
		if err != nil {
			log.Printf("GBTT departure time parse error: %v", err)
			valid = false
		}

		realtimeDeparture, err := trainTimeToGoTime(date, s.LocationDetail.RealtimeDeparture)
		if err != nil {
			log.Printf("Realtime departure time parse error: %v", err)
			valid = false
		}
		
		lateness := realtimeDeparture.Sub(gbttDeparture)
		// If we're ludicrously early, it's more likely that we're slightly late
		if lateness <= -12 * time.Hour {
			realtimeDeparture = realtimeDeparture.Add(24 * time.Hour)
			lateness = lateness + 12 * time.Hour
		}

		return WrappedService{
			Valid:             valid,
			RequestDate:       date,
			GBTTDeparture:     gbttDeparture,
			RealtimeDeparture: realtimeDeparture,
			Lateness:          realtimeDeparture.Sub(gbttDeparture),
			S:                 s,
		}
	}))
}
