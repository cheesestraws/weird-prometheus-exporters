package truenas

import (
	"context"
	"os"
	"testing"
	
	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

func TestBasicGet(t *testing.T) {
	baseurl := os.Getenv("NASURL")
	user := os.Getenv("NASUSER")
	pass := os.Getenv("NASPASS")

	if baseurl == "" || user == "" || pass == "" {
		t.Skip("set NASURL, NASUSER and NASPASS to the details of your TrueNAS to run integration tests")
	}

	c := NewClient(baseurl, user, pass, nil)

	s, err := BasicGet[string](c, context.Background(), "/core/ping")
	if err != nil {
		t.Errorf("ping err: %v", err)
	}

	if s != "pong" {
		t.Errorf("expected 'pong', got %q", s)
	}
}

func TestAlertList(t *testing.T) {
	baseurl := os.Getenv("NASURL")
	user := os.Getenv("NASUSER")
	pass := os.Getenv("NASPASS")

	if baseurl == "" || user == "" || pass == "" {
		t.Skip("set NASURL, NASUSER and NASPASS to the details of your TrueNAS to run integration tests")
	}

	c := NewClient(baseurl, user, pass, nil)

	as, err := c.AlertList(context.Background())
	if err != nil {
		t.Errorf("alert err: %v", err)
	}

	dismissed := fn.Count(as, func(a Alert) bool {
		return a.Dismissed
	})
	
	t.Logf("%d alerts, %d dismissed", len(as), dismissed)
}

func TestPools(t *testing.T) {
	baseurl := os.Getenv("NASURL")
	user := os.Getenv("NASUSER")
	pass := os.Getenv("NASPASS")

	if baseurl == "" || user == "" || pass == "" {
		t.Skip("set NASURL, NASUSER and NASPASS to the details of your TrueNAS to run integration tests")
	}

	c := NewClient(baseurl, user, pass, nil)

	ps, err := c.Pools(context.Background())
	if err != nil {
		t.Errorf("pools err: %v", err)
	}

	t.Logf("%d pools", len(ps))
}

func TestCloudSync(t *testing.T) {
	baseurl := os.Getenv("NASURL")
	user := os.Getenv("NASUSER")
	pass := os.Getenv("NASPASS")

	if baseurl == "" || user == "" || pass == "" {
		t.Skip("set NASURL, NASUSER and NASPASS to the details of your TrueNAS to run integration tests")
	}

	c := NewClient(baseurl, user, pass, nil)

	cs, err := c.CloudSyncs(context.Background())
	if err != nil {
		t.Errorf("cloudsyncs err: %v", err)
	}

	t.Logf("%+v", cs)
}