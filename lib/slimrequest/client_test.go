package slimrequest

import (
	"context"
	"os"
	"testing"
)

func TestBasicDo(t *testing.T) {
	baseurl := os.Getenv("LMSURL")
	if baseurl == "" {
		t.Skip("set LMSURL to the URL of your LMS instance to run integration tests")
	}

	c := NewClient(baseurl, nil)

	r := NewRequest("0", []string{"serverstatus"})

	resp, err := Do[ServerStatus](c, context.Background(), r)
	if err != nil {
		t.Errorf("Do returned error: %v", err)
	}

	t.Logf("please verify this matches your setup for a serverstatus command: %+v", resp)
}

func TestServerStatusExtended(t *testing.T) {
	baseurl := os.Getenv("LMSURL")
	if baseurl == "" {
		t.Skip("set LMSURL to the URL of your LMS instance to run integration tests")
	}

	c := NewClient(baseurl, nil)
	
	stat, err := c.ExtendedServerStatus(context.Background())
	if err != nil {
		t.Errorf("ExtendedServerStatus returned error: %v", err)
	}
	
	t.Logf("please verify this matches your setup : %+v", stat)
}

func TestPlayerStatus(t *testing.T) {
	baseurl := os.Getenv("LMSURL")
	if baseurl == "" {
		t.Skip("set LMSURL to the URL of your LMS instance to run integration tests")
	}

	c := NewClient(baseurl, nil)
	
	stat, err := c.ExtendedServerStatus(context.Background())
	if err != nil {
		t.Errorf("ExtendedServerStatus returned error: %v", err)
	}
	
	id := stat.Players[0].PlayerID
	
	sstat, err := c.PlayerStatus(context.Background(), id)
	if err != nil {
		t.Errorf("PlayerStatus returned error: %v", err)
	}
	
	t.Logf("please verify this matches your setup : %+v", sstat)
}