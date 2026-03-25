package rflookup

import (
	"testing"
)

func TestHostnames(t *testing.T) {
	// Meh
	centreFrequency := 433920000

	type test struct {
		labels       map[string]string
		expectedFQDN string
	}

	tests := []test{
		test{
			labels: map[string]string{
				"channel": "1",
				"id":      "219",
				"model":   "LaCrosse-TX141THBv2",
			},
			expectedFQDN: "lacrosse-tx141thbv2-channel-1-id-219.433.rf",
		},

		test{
			labels: map[string]string{
				"data":  "212de30",
				"model": "Doorbell_World_EV1527",
			},
			expectedFQDN: "doorbell_world_ev1527-data-212de30.433.rf",
		},
	}

	for _, test := range tests {
		fqdn, err := MkLookupHostname(test.labels, centreFrequency)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		if fqdn != test.expectedFQDN {
			t.Errorf("expected %s, got %s", test.expectedFQDN, fqdn)
		}
	}

}

var enableTestsThatOnlyWorkOnCheeseysNetwork bool = false

func TestResolution(t *testing.T) {
	if enableTestsThatOnlyWorkOnCheeseysNetwork {
		centreFrequency := 433920000
		labels := map[string]string{
			"data":  "test",
			"model": "Doorbell_World_EV1527",
		}

		hn, err := Lookup(labels, centreFrequency)

		if err != nil {
			t.Fatalf("err: %v", err)
		}
		
		if hn != "test-doorbell" {
			t.Fatalf("got wrong hostname: %s", hn)
		}
	}
}
