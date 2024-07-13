package main

import (
	"time"
	
	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

type TimeWindow struct {
	Name string
	From func() time.Time
	To   func() time.Time
}

type TimeWindowSnapshot struct {
	Name string
	From time.Time
	To time.Time
}

func (t TimeWindow) Snapshot() TimeWindowSnapshot {
	return TimeWindowSnapshot{
		Name: t.Name,
		From: t.From(),
		To: t.To(),
	}
}

type TimeWindows []TimeWindow

func (t TimeWindows) Snapshot() []TimeWindowSnapshot {
	return fn.Map(t, func(w TimeWindow) TimeWindowSnapshot {
		return w.Snapshot()
	})
}

type Config struct {
	TimeWindows TimeWindows
}
