package main

import (
	"bytes"
	"fmt"
	"reflect"
	"time"
)

type Summaries struct {
	Window        string
	StationName   string
	StationCRS    string
	StationTIPLOC string

	Origin      *Location
	Destination *Location

	NumTrains          int `prometheus:"num_trains"`
	NumLateTrains      int `prometheus:"num_late_trains"`
	NumCancelledTrains int `prometheus:"num_cancelled_trains"`

	AvgLateTime   time.Duration `prometheus:"avg_late_time"`
	WorstLateTime time.Duration `prometheus:"worst_late_time"`

	CancelReasons map[string]int `prometheus_map:"cancel_reasons" prometheus_map_key:"reason"`
}

func (ss WrappedServices) Summarise(window string) *Summaries {
	sum := &Summaries{
		Window:        window,
		CancelReasons: make(map[string]int),
	}

	var latenessAccumulator time.Duration

	for _, s := range ss {
		if !s.Valid {
			continue
		}

		sum.NumTrains++

		if s.Lateness > 5*time.Minute {
			sum.NumLateTrains++
		}

		if s.S.LocationDetail.DisplayAs == "CANCELLED_CALL" {
			sum.NumCancelledTrains++
		}

		latenessAccumulator += s.Lateness

		if s.Lateness > sum.WorstLateTime {
			sum.WorstLateTime = s.Lateness
		}

		cc :=  s.S.LocationDetail.CancelReasonCode
		if s.S.LocationDetail.CancelReasonShortText != "" {
			sum.CancelReasons[cc + " - " + s.S.LocationDetail.CancelReasonShortText]++
		} else if s.S.LocationDetail.CancelReasonLongText != "" {
			sum.CancelReasons[cc + " - " + s.S.LocationDetail.CancelReasonLongText]++
		} else if s.S.LocationDetail.CancelReasonCode != "" {
			sum.CancelReasons[s.S.LocationDetail.CancelReasonCode]++
		}

		sum.StationName = s.S.LocationDetail.Description
		sum.StationCRS = s.S.LocationDetail.CRS
		sum.StationTIPLOC = s.S.LocationDetail.TIPLOC
	}

	if sum.NumTrains != 0 {
		sum.AvgLateTime = latenessAccumulator / time.Duration(sum.NumTrains)
	}

	return sum
}

func (s *Summaries) metricName(baseMetricName string) string {
	mn := fmt.Sprintf("realtimetrains_%s", baseMetricName)

	if s.Origin != nil {
		mn += "_from"
	}

	if s.Destination != nil {
		mn += "_to"
	}
	
	return mn
}

func (s *Summaries) metricLabels() string {
	ll := fmt.Sprintf("timewindow=%q,station=%q,crs=%q,tiploc=%q", s.Window, s.StationName, s.StationCRS, s.StationTIPLOC)

	if s.Origin != nil {
		ll += ","
		ll += s.Origin.PrometheusLabels("from_")
	}

	if s.Destination != nil {
		ll += ","
		ll += s.Destination.PrometheusLabels("to_")
	}

	return ll
}

func (s *Summaries) Prometheise() []byte {
	labels := s.metricLabels()
	bs := &bytes.Buffer{}

	mtype := func(metric string, t string) {
		fmt.Fprintf(bs, "# TYPE %s %s\n", metric, t)
	}

	val := func(metric string, extraLabels string, i interface{}) {
		var valstr string
		switch v := i.(type) {
		case time.Duration:
			valstr = fmt.Sprintf("%v", v.Seconds())
		default:
			valstr = fmt.Sprintf("%v", v)
		}
		fmt.Fprintf(bs, "%s{%s%s} %s\n", metric, labels, extraLabels, valstr)
	}

	// Now iterate through the fields
	theStruct := reflect.ValueOf(*s)
	myType := theStruct.Type()

	for i := 0; i < myType.NumField(); i++ {
		fld := myType.Field(i)
		// is it a scalar?
		basename := fld.Tag.Get("prometheus")
		if basename != "" {
			realName := s.metricName(basename)
			value := theStruct.Field(i).Interface()
			mtype(realName, "gauge")
			val(realName, "", value)
		}

		// is it a map?
		mapname := fld.Tag.Get("prometheus_map")
		keyname := fld.Tag.Get("prometheus_map_key")
		if mapname != "" && keyname != "" && fld.Type.Kind() == reflect.Map {
			// Iterate over it
			realName := s.metricName(mapname)

			mtype(realName, "gauge")
			iter := theStruct.Field(i).MapRange()
			for iter.Next() {
				tagValue := fmt.Sprintf("%v", iter.Key().Interface())
				metricValue := iter.Value().Interface()
				extraTags := fmt.Sprintf(",%s=%q", keyname, tagValue)
				val(realName, extraTags, metricValue)
			}
		}
	}

	return bs.Bytes()
}
