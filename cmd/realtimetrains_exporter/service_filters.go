package main

import (
	"time"
	
	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

func (ss WrappedServices) ByTimeWindow(from time.Time, to time.Time) WrappedServices {
	return fn.Filter(ss, func(s WrappedService) bool {
		return (s.GBTTDeparture.After(from) || s.RealtimeDeparture.After(from)) &&
			(s.GBTTDeparture.Before(to) || s.RealtimeDeparture.After(to))
	})
}
