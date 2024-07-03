package main

/* ea.go contains types that correspond to json chunks in the
   realtimeish river level dataset */

import (
	"encoding/json"
	"fmt"
	"time"
)

// Deal with broken timestamps that are missing time zone
type timestamp time.Time

func (t *timestamp) UnmarshalJSON(bs []byte) error {
	// This is a bodge.  Thank goodness we're hangin' around near
	// the Greenwich meridian.
	var timeStr string
	err1 := json.Unmarshal(bs, &timeStr)
	if err1 != nil {
		return err1
	}

	tm, err1 := time.Parse(time.RFC3339, timeStr)
	if err1 == nil {
		*t = timestamp(tm)
		return nil
	}

	// Let's try parsing it with a forced zero time zone
	timeStr = timeStr + "Z"
	tm, err2 := time.Parse(time.RFC3339, timeStr)
	if err2 == nil {
		*t = timestamp(tm)
		return nil
	}

	return err1
}

func (t *timestamp) String() string {
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("%+v", *t)
}

type station struct {
	Items stationItems `json:"items"`
}

type stationItems struct {
	Label      string   `json:"label"`
	Notation   string   `json:"notation"`
	RiverName  string   `json:"riverName"`
	StageScale *scale   `json:"stageScale"`
	Measures   measures `json:"measures"`
}

type scale struct {
	HighestRecent    reading  `json:"highestRecent"`
	MaxOnRecord      reading  `json:"maxOnRecord"`
	MinOnRecord      reading  `json:"minOnRecord"`
	ScaleMax         float64  `json:"scaleMax"`
	TypicalRangeHigh *float64 `json:"typicalRangeHigh"`
	TypicalRangeLow  *float64 `json:"typicalRangeLow"`
}

type readingFields struct {
	Valid     bool
	Timestamp *timestamp `json:"dateTime"`
	Value     float64    `json:"value"`
}

type reading struct {
	readingFields
}

func (r *reading) UnmarshalJSON(bs []byte) error {
	// Sometimes we just have an ID instead of a reading.  Why?  Dunno.
	var dummyStr string
	err := json.Unmarshal(bs, &dummyStr)
	if err == nil {
		// yeah we just got an ID
		*r = reading{readingFields{
			Valid: false,
		}}
		return nil
	}

	// not a string, let's be responsible
	err = json.Unmarshal(bs, &r.readingFields)
	if err != nil {
		return err
	}

	r.readingFields.Valid = true
	return nil
}

func (r *reading) String() string {
	if r == nil {
		return "nil"
	}
	return fmt.Sprintf("%+v", r.readingFields)
}

type measure struct {
	Parameter     string  `json:"parameter"`
	Qualifier     string  `json:"qualifier"`
	LatestReading reading `json:"latestReading"`
	UnitName      string  `json:"unitName"`
}

type measures []measure

func (ms *measures) UnmarshalJSON(bs []byte) error {
	// We might instead just have a single measure, do we?
	var m measure
	err := json.Unmarshal(bs, &m)
	if err == nil {
		// yeah we have been given only a single measure
		// and that not wrapped in an array.
		// Ho hum.

		*ms = measures{m}
		return nil
	}

	var mss []measure
	err = json.Unmarshal(bs, &mss)
	if err != nil {
		return err
	}

	*ms = mss
	return nil
}

func (ms *measures) String() string {
	if ms == nil {
		return "nil"
	}

	return fmt.Sprintf("%+v", *ms)
}
