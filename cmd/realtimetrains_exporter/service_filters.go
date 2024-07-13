package main

import (
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"

	rtt "github.com/cheesestraws/gortt"
)

func (ss WrappedServices) ByTimeWindow(from time.Time, to time.Time) WrappedServices {
	return fn.Filter(ss, func(s WrappedService) bool {
		return (s.GBTTDeparture.After(from) || s.RealtimeDeparture.After(from)) &&
			(s.GBTTDeparture.Before(to) || s.RealtimeDeparture.Before(to))
	})
}

func (ss WrappedServices) ByDestinationTIPLOC(tiploc string) WrappedServices {
	return fn.Filter(ss, func(s WrappedService) bool {
		return fn.Contains(s.S.LocationDetail.Destination, func(p rtt.RTTPair) bool {
			return p.TIPLOC == tiploc
		})
	})
}

func (ss WrappedServices) ByOriginTIPLOC(tiploc string) WrappedServices {
	return fn.Filter(ss, func(s WrappedService) bool {
		return fn.Contains(s.S.LocationDetail.Origin, func(p rtt.RTTPair) bool {
			return p.TIPLOC == tiploc
		})
	})
}
