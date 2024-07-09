package netatalk

import (
	"fmt"
	"os"
	"testing"
)

func TestCLIQuerier(t *testing.T) {
	if os.Getenv("WITHAT") == "" {
		t.Skipf("set WITHAT to run tests that require netatalk")
	}

	var q cliQuerier

	t.Run("GetZones", func(t *testing.T) {
		zs, err := q.GetZones()
		if err != nil {
			t.Errorf("getzones returned err: %v", err)
			return
		}

		t.Logf("Please check zones against real network configuration.")
		t.Logf("Zones: %+v", zs)
	})

	t.Run("NBPLookup", func(t *testing.T) {
		zs, err := q.GetZones()
		if err != nil {
			t.Errorf("getzones returned err: %v", err)
			return
		}

		for _, zone := range zs {
			t.Logf(" ")
			t.Logf("Zone: %s", zone)
			m, err := q.NBPLookup(fmt.Sprintf("@%s", zone))
			if err != nil {
				t.Errorf("NBPLookup returned err: %v", err)
				return
			}
	
			for k, v := range m {
				t.Logf("    %s => %s", k, v)
			}
			
			m, err = q.NBPLookup(fmt.Sprintf("=:AppleRouter@%s", zone))
			if err != nil {
				t.Errorf("NBPLookup returned err: %v", err)
				return
			}

			if len(m) > 1 {
				t.Logf("    (zone contains an Apple Router)")
			}

		}
	})

}
