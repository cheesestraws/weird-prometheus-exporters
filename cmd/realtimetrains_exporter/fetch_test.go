package main

import (
	"reflect"
	"testing"
	"time"
)

const dtformat = "2006-01-02 15:04:05"

func mkdate(s string) time.Time {
	d, _ := time.Parse("2006-01-02", s)
	d = d.Local()
	return d
}

func mkdatetime(s string) *time.Time {
	d, _ := time.Parse(dtformat, s)
	d = d.Local()
	return &d
}

func TestMakeFetches(t *testing.T) {
	t.Run("overlappage", func(t *testing.T) {
		from, _ := time.Parse(dtformat, "2003-02-01 03:04:05")
		to := from.Add(-3 * time.Hour)
		f := MakeFetches("SAL", from, to)
		if f != nil {
			t.Errorf("unexpectedly not nil when from > to")
		}
	})

	t.Run("same day", func(t *testing.T) {
		from, _ := time.Parse(dtformat, "2003-02-01 03:04:05")
		to := from.Add(3 * time.Hour)
		f := MakeFetches("SAL", from, to)

		expected := Fetches{
			{
				Station: "SAL",
				Date:    mkdate("2003-02-01"),
				From:    mkdatetime("2003-02-01 03:04:05"),
				To:      mkdatetime("2003-02-01 06:04:05"),
			},
		}

		if !reflect.DeepEqual(expected, f) {
			t.Errorf("mismatch expected vs result: %+v vs %+v", expected, f)
		}

	})

	t.Run("two days", func(t *testing.T) {
		from, _ := time.Parse(dtformat, "2003-02-01 03:04:05")
		to := from.Add(26 * time.Hour)
		f := MakeFetches("SAL", from, to)

		expected := Fetches{
			{
				Station: "SAL",
				Date:    mkdate("2003-02-01"),
				From:    mkdatetime("2003-02-01 03:04:05"),
				To:      nil,
			},
			{
				Station: "SAL",
				Date:    mkdate("2003-02-02"),
				From:    nil,
				To:      mkdatetime("2003-02-02 05:04:05"),
			},
		}

		if !reflect.DeepEqual(expected, f) {
			t.Errorf("mismatch expected vs result: %+v vs %+v", expected, f)
		}
	})

}
