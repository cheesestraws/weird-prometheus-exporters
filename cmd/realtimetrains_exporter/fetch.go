package main

import (
	"context"
	"time"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

type Fetch struct {
	Station string
	Date    time.Time
	From    *time.Time
	To      *time.Time
}

type FetchResult struct {
	Fetch
	//Results *RTTLocationLineup
}

func (f *Fetch) Do(ctx context.Context, rttapi interface{}) (*FetchResult, error) {
	return nil, nil
}

// truncdate, truncates to the nearest local day.  it's a pun, geddit
func truncdate(t time.Time) time.Time {
	yyyy, mm, dd := t.Date()

	return time.Date(yyyy, mm, dd, 0, 0, 0, 0, time.Local)
}

type Fetches []Fetch

func MakeFetches(station string, from time.Time, to time.Time) Fetches {
	if to.Before(from) {
		return nil
	}

	// programmers are bad at dates
	// i am no exception

	var dates []time.Time

	iterate := truncdate(from)
	for iterate.Before(to) {
		dates = append(dates, iterate)
		// 26 because midnight + 26 is always the next day even if
		// summer time, I hate this
		iterate = truncdate(iterate.Add(26 * time.Hour))
	}

	fetches := fn.Map(dates, func(d time.Time) Fetch {
		return Fetch{
			Station: station,
			Date:    d,
		}
	})

	f := from.Local()
	fetches[0].From = &f

	t := to.Local()
	fetches[len(fetches)-1].To = &t

	return fetches
}
