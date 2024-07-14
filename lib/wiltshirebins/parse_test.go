package wiltshirebins

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	rdr := strings.NewReader(testCalendar)

	var days Calendar

	days, err := parse(rdr)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if days[0].HouseholdWaste != 0 || days[0].Recycling != 0 {
		t.Errorf("1st should have no collection")
	}

	if days[4].HouseholdWaste != 0 || days[4].Recycling != 1 {
		t.Errorf("5th should have recycling")
	}

	if days[11].HouseholdWaste != 1 || days[11].Recycling != 0 {
		t.Errorf("12th should have household waste")
	}
}
