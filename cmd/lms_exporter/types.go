package main

import (
	"strings"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/slimrequest"
)

type Status struct {
	Server       slimrequest.ExtendedServerStatus
	PlayerInfo   map[string]slimrequest.PlayerDetails
	PlayerStatus map[string]slimrequest.PlayerStatus
}

type PlayerIdentifier struct {
	Name      string `prometheus_label:"player_name"`
	ID        string `prometheus_label:"player_id"`
	ModelName string `prometheus_label:"player_model"`
}

type Summary struct {
	TotalSongs    int     `prometheus:"total_songs"`
	TotalArtists  int     `prometheus:"total_artists"`
	TotalAlbums   int     `prometheus:"total_albums"`
	TotalGenres   int     `prometheus:"total_genres"`
	TotalDuration float64 `prometheus:"total_duration" prometheus_help:"in seconds"`

	PlayerCount int `prometheus:"player_count"`

	PlayerConnected map[PlayerIdentifier]int `prometheus_map:"player_connected"`
	PlayerOn        map[PlayerIdentifier]int `prometheus_map:"player_on"`
	Volume          map[PlayerIdentifier]int `prometheus_map:"player_volume"`
	SignalStrength  map[PlayerIdentifier]int `prometheus_map:"player_signalstrength"`

	Models         map[string]int `prometheus_map:"players_by_model" prometheus_map_key:"player_model"`
	Modes          map[string]int `prometheus_map:"players_by_mode" prometheus_map_key:"mode"`
	SyncGroupSizes map[string]int `prometheus_map:"sync_groups" prometheus_map_key:"master"`
}

func (s *Status) Summarise() Summary {
	sum := Summary{
		TotalSongs:    s.Server.TotalSongs,
		TotalArtists:  s.Server.TotalArtists,
		TotalAlbums:   s.Server.TotalAlbums,
		TotalGenres:   s.Server.TotalGenres,
		TotalDuration: s.Server.TotalDuration,

		PlayerCount: s.Server.PlayerCount,

		PlayerConnected: make(map[PlayerIdentifier]int),
		PlayerOn:        make(map[PlayerIdentifier]int),
		Volume:          make(map[PlayerIdentifier]int),
		SignalStrength:  make(map[PlayerIdentifier]int),

		Models:         make(map[string]int),
		Modes:          map[string]int{"play": 0, "pause": 0, "stop": 0},
		SyncGroupSizes: make(map[string]int),
	}

	for id, info := range s.PlayerInfo {
		idx := PlayerIdentifier{
			Name:      info.Name,
			ID:        info.PlayerID,
			ModelName: info.ModelName,
		}

		sum.PlayerConnected[idx] = info.Connected
		sum.PlayerOn[idx] = info.Power

		sum.Volume[idx] = s.PlayerStatus[id].MixerVolume
		sum.SignalStrength[idx] = s.PlayerStatus[id].SignalStrength

		sum.Models[info.ModelName]++
		sum.Modes[s.PlayerStatus[id].Mode]++
	}

	// Now deal with sync groups.  All are 0 sized to start with:
	for _, st := range s.PlayerStatus {
		sum.SyncGroupSizes[st.PlayerName] = 0
	}

	// Then fill in those that aren't 0
	for _, st := range s.PlayerStatus {
		if st.SyncMaster == "" && st.SyncSlaves == "" {
			continue
		}

		// Find the master
		m, ok := s.PlayerStatus[st.SyncMaster]
		var mn string
		if ok {
			mn = m.PlayerName
		} else {
			mn = st.SyncMaster
		}

		// And the rest
		ss := strings.Split(st.SyncSlaves, ",")
		c := len(ss)

		sum.SyncGroupSizes[mn] = c
	}

	return sum
}
