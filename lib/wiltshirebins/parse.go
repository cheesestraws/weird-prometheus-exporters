package wiltshirebins

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func parse(rdr io.Reader) (Calendar, error) {
	var days Calendar

	currentDay := 0
	z := html.NewTokenizer(rdr)
outer:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err != io.EOF {
				return days, err
			}
			break outer

		case html.TextToken:
			s := string(bytes.TrimSpace(z.Text()))
			if len(s) != 0 {
				// is this a date?
				num, err := strconv.Atoi(s)
				if err == nil {
					currentDay = num
					continue outer
				}

				// Is it a recycling day?
				if strings.HasPrefix(s, "Mixed") {
					days[currentDay-1].Recycling = 1
				}

				if strings.HasPrefix(s, "Household") {
					days[currentDay-1].HouseholdWaste = 1
				}
			}

		case html.StartTagToken:
			// If we get to the footer, stop
			var k, v []byte

			tn, hasAttr := z.TagName()
			if string(tn) == "div" {
				for hasAttr {
					k, v, hasAttr = z.TagAttr()
					if string(k) == "class" && strings.Contains(string(v), "-foot-") {
						break outer
					}
				}
			}
		}
	}

	// perhaps some extra sanity checks should go here when I'm saner

	return days, nil
}
