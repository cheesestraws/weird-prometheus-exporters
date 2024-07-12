package main

import (
	"reflect"
	"testing"
	"time"

	rtt "github.com/cheesestraws/gortt"
)

func TestLocationLineupToServices(t *testing.T) {
	sdate := func(s string) time.Time {
		d, _ := time.Parse("2006-01-02", s)
		return d
	}
	
	sdatetime := func(s string) time.Time {
		d, _ := time.Parse("2006-01-02 15:04:05", s)
		return d
	}


	d := sdate("2023-04-05")
	ll := rtt.RTTLocationLineup{
		Services: []rtt.RTTLocationContainer{
			{
				ServiceUID: "one",
				LocationDetail: rtt.RTTLocation{
					GBTTBookedDeparture: "0115",
					RealtimeDeparture:   "0120",
				},
			},
			{
				ServiceUID: "two",
				LocationDetail: rtt.RTTLocation{
					GBTTBookedDeparture: "0215",
					RealtimeDeparture:   "0210",
				},
			},
			{
				ServiceUID: "three",
				LocationDetail: rtt.RTTLocation{
					GBTTBookedDeparture: "2330",
					RealtimeDeparture:   "0015",
				},
			},

		},
	}

	ws := LocationLineupToServices(ll, d)
	expected := WrappedServices{
		{
			Valid:       true,
			RequestDate: d,
			GBTTDeparture: sdatetime("2023-04-05 01:15:00"),
			RealtimeDeparture: sdatetime("2023-04-05 01:20:00"),
			Lateness: 5 * time.Minute,

			
			S: rtt.RTTLocationContainer{
				ServiceUID: "one",
				LocationDetail: rtt.RTTLocation{
					GBTTBookedDeparture: "0115",
					RealtimeDeparture:   "0120",
				},
			},
		},
		{
			Valid:       true,
			RequestDate: d,
			GBTTDeparture: sdatetime("2023-04-05 02:15:00"),
			RealtimeDeparture: sdatetime("2023-04-05 02:10:00"),
			Lateness: -5 * time.Minute,

			
			S: rtt.RTTLocationContainer{
				ServiceUID: "two",
				LocationDetail: rtt.RTTLocation{
					GBTTBookedDeparture: "0215",
					RealtimeDeparture:   "0210",
				},
			},
		},
		{
			Valid:       true,
			RequestDate: d,
			GBTTDeparture: sdatetime("2023-04-05 23:30:00"),
			RealtimeDeparture: sdatetime("2023-04-06 00:15:00"),
			Lateness: 45 * time.Minute,

			
			S: rtt.RTTLocationContainer{
				ServiceUID: "three",
				LocationDetail: rtt.RTTLocation{
					GBTTBookedDeparture: "2330",
					RealtimeDeparture:   "0015",
				},
			},
		},

	}

	if !reflect.DeepEqual(expected, ws) {
		t.Errorf("bad result: expected \n %+v\n, got\n %+v\n", expected, ws)
	}
}
