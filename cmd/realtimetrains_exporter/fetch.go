package main

import (
	"context"
	"errors"
	"time"

	rtt "github.com/cheesestraws/gortt"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

type Fetch struct {
	Station string
	Date    time.Time
	From    *time.Time
	To      *time.Time
}

func (f *Fetch) Do(ctx context.Context, cli *rtt.Client) (rtt.RTTLocationDetail, WrappedServices, error) {
	ll, err := cli.GetServicesFromStationForDate(ctx, f.Station, f.Date)
	if err != nil {
		return rtt.RTTLocationDetail{}, nil, err
	}

	if ll == nil {
		return rtt.RTTLocationDetail{}, nil, errors.New("silently got nil lineup; probably a bug")
	}

	return ll.Location, LocationLineupToServices(*ll, f.Date), nil
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

func (fs Fetches) Do(ctx context.Context, cli *rtt.Client) (rtt.RTTLocationDetail, WrappedServices, error) {
	var ss WrappedServices
	var ld rtt.RTTLocationDetail
	for _, f := range fs {
		ld, ws, err := f.Do(ctx, cli)
		if err != nil {
			return ld, ss, err
		}

		ss = append(ss, ws...)
	}

	return ld, ss, nil
}
