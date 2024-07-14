package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

type Fudge struct {
	TypicalLow  float64
	TypicalHigh float64
}

func ParseFudge(s string) (string, Fudge, error) {
	var f Fudge
	var err error

	ss := strings.Split(s, ":")
	if len(ss) != 3 {
		return "", f, fmt.Errorf("fudge %q has wrong number of elements (supply station ID, low and high)", s)
	}

	station := ss[0]

	f.TypicalLow, err = strconv.ParseFloat(ss[1], 64)
	if err != nil {
		return "", f, fmt.Errorf("fudge for station %q: low %q is not a number", station, ss[1])
	}

	f.TypicalHigh, err = strconv.ParseFloat(ss[2], 64)
	if err != nil {
		return "", f, fmt.Errorf("fudge for station %q: high %q is not a number", station, ss[2])
	}

	return station, f, nil
}

type Fudges map[string]Fudge

func ParseFudges(s string) (Fudges, error) {
	if len(s) == 0 {
		return nil, nil
	}

	ss := strings.Split(s, ",")
	fm, err := fn.Errmapmap(ss, ParseFudge)
	return Fudges(fm), err
}

func (f Fudges) Get(station string) (low float64, high float64, ok bool) {
	if f == nil {
		return 0, 0, false
	}

	fudge, ok := f[station]
	if !ok {
		return 0, 0, false
	}

	return fudge.TypicalLow, fudge.TypicalHigh, true
}
