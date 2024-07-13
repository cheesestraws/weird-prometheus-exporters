package main

import (
	"time"
)

// this should be a config file but I'm too tired

var config = Config{
	TimeWindows: []TimeWindow{
		{
			Name: "1h back, 30m ahead",
			From: func() time.Time { return time.Now().Add(-1 * time.Hour) },
			To:   func() time.Time { return time.Now().Add(30 * time.Minute) },
		},
		
		{
			Name: "2h ahead",
			From: func() time.Time { return time.Now() },
			To:   func() time.Time { return time.Now().Add(2 * time.Hour) },
		},
	},
}
