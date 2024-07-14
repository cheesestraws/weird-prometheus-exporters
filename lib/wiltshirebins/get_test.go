package wiltshirebins

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	// Really run tests that talk to the 'net?
	if _, ok := os.LookupEnv("MISTERT"); !ok {
		t.Skip("set MISTERT in the environment to run integration tests")
	}

	now := time.Now()
	month := int(now.Month())
	year := int(now.Year())

	// The addresses used in this integration test have NOTHING TO DO WITH ME.
	// Do not knock on their doors and say you don't like my code, they won't
	// know what you're talking about.  Raise an abusive issue on github instead
	// as is tradition.

	t.Run("residential, should get collections", func(t *testing.T) {
		days, err := DefaultClient.Get(context.Background(), month, year, "SP1 2EJ", "010010443460")
		if err != nil {
			t.Errorf("Get() errored: %v", err)
		}

		var totalHW int
		var totalRe int
		var totalErrors int
		for _, v := range days {
			totalHW += v.HouseholdWaste
			totalRe += v.Recycling
			totalErrors += len(v.Errors)
		}

		if totalHW < 1 || totalRe < 1 || totalErrors > 0 {
			t.Errorf("bad data for residential?")
		}

	})

	t.Run("nonresidential, should not get collections", func(t *testing.T) {
		days, err := DefaultClient.Get(context.Background(), month, year, "SP1 2EJ", "010093279572")
		if err != nil {
			t.Errorf("Get() errored: %v", err)
		}

		var totalHW int
		var totalRe int
		var totalErrors int
		for _, v := range days {
			totalHW += v.HouseholdWaste
			totalRe += v.Recycling
			totalErrors += len(v.Errors)
		}

		// errors are >= 31 because sanity checking should pick up on
		// data being implausible
		if totalHW > 0 || totalRe > 0 || totalErrors < 31 {
			t.Errorf("bad data for nonresidential?")
		}

	})

}
