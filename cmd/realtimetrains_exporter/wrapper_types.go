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

	StationName   string
	StationCRS    string
	StationTIPLOC string

	RequestDate time.Time

	GBTTDeparture     time.Time
	RealtimeDeparture time.Time
	Lateness          time.Duration

	Cancelled bool

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

		// do we have a cancellation reason?
		cancelled := false
		if s.LocationDetail.CancelReasonShortText != "" {
			cancelled = true
		} else if s.LocationDetail.CancelReasonLongText != "" {
			cancelled = true
		} else if s.LocationDetail.CancelReasonCode != "" {
			cancelled = true
		} else if s.LocationDetail.DisplayAs == "CANCELLED_CALL" {
			cancelled = true
		}
		gbttDeparture, err := trainTimeToGoTime(date, s.LocationDetail.GBTTBookedDeparture)
		if err != nil {
			log.Printf("GBTT departure time parse error: %v", err)
			valid = false
		}

		realtimeDeparture, err := trainTimeToGoTime(date, s.LocationDetail.RealtimeDeparture)
		if err != nil && s.ServiceType != "bus" {
			log.Printf("Realtime departure time parse error: %v (GBTT time is %v)", err, gbttDeparture)
			valid = false
		} else if s.ServiceType == "bus" {
			// whistle while you bodge
			realtimeDeparture = gbttDeparture
		}

		if cancelled {
			// real time for cancelled trains is 'never'.
			// Bodge it so they do something useful
			realtimeDeparture = gbttDeparture
		}

		lateness := realtimeDeparture.Sub(gbttDeparture)
		// If we're ludicrously early, it's more likely that we're slightly late
		if lateness <= -12*time.Hour {
			realtimeDeparture = realtimeDeparture.Add(24 * time.Hour)
			lateness = lateness + 12*time.Hour
		}

		return WrappedService{
			Valid: valid,

			StationName:   ll.Location.Name,
			StationCRS:    ll.Location.CRS,
			StationTIPLOC: ll.Location.TIPLOC,

			RequestDate:       date,
			GBTTDeparture:     gbttDeparture,
			Cancelled:         cancelled,
			RealtimeDeparture: realtimeDeparture,
			Lateness:          realtimeDeparture.Sub(gbttDeparture),
			S:                 s,
		}
	}))
}

type Location struct {
	Name   string
	TIPLOC string
}

func (l Location) PrometheusLabels(prefix string) string {
	return fmt.Sprintf("%s%s=%q,%s%s=%q",
		prefix, "name", l.Name,
		prefix, "tiploc", l.TIPLOC)
}

func (w WrappedService) Origins() []Location {
	accum := make(map[Location]struct{})

	for _, v := range w.S.LocationDetail.Origin {
		loc := Location{
			Name:   v.Description,
			TIPLOC: v.TIPLOC,
		}

		accum[loc] = struct{}{}
	}

	return fn.Mapkeymap(accum, fn.Id[Location])
}

func (w WrappedService) Destinations() []Location {
	accum := make(map[Location]struct{})

	for _, v := range w.S.LocationDetail.Destination {
		loc := Location{
			Name:   v.Description,
			TIPLOC: v.TIPLOC,
		}

		accum[loc] = struct{}{}
	}

	return fn.Mapkeymap(accum, fn.Id[Location])
}

func (w WrappedServices) Origins() []Location {
	var locations []Location

	for _, v := range w {
		locations = append(locations, v.Origins()...)
	}

	return fn.Dedupe(locations)
}

func (w WrappedServices) Destinations() []Location {
	var locations []Location

	for _, v := range w {
		locations = append(locations, v.Destinations()...)
	}

	return fn.Dedupe(locations)
}
