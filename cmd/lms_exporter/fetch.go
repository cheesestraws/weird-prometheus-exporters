package main

import (
	"context"
	"fmt"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/slimrequest"
)

func fetchStatus(ctx context.Context, baseURL string) (Status, error) {
	var s Status

	c := slimrequest.NewClient(baseURL, nil)

	// fetch server status
	extstat, err := c.ExtendedServerStatus(ctx)
	if err != nil {
		return s, fmt.Errorf("ExtendedServerStatus returned error: %v", err)
	}

	s.Server = extstat

	// copy player info into playerinfo map
	s.PlayerInfo = make(map[string]slimrequest.PlayerDetails)
	for _, v := range extstat.Players {
		s.PlayerInfo[v.PlayerID] = v
	}

	// look up player status
	s.PlayerStatus = make(map[string]slimrequest.PlayerStatus)
	for k := range s.PlayerInfo {
		p, err := c.PlayerStatus(ctx, k)
		if err != nil {
			return s, fmt.Errorf("PlayerStatus (for player %q) returned error: %v", k, err)
		}
		s.PlayerStatus[k] = p
	}

	return s, nil
}
