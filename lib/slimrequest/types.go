package slimrequest

type ServerStatus struct {
	IP          string `json:"ip"`
	LastScan    string `json:"lastscan"`
	Version     string `json:"version"`
	MAC         string `json:"mac"`
	HTTPPort    string `json:"httpport"`
	PlayerCount int    `json:"player count"`
	UUID        string `json:"uuid"`

	TotalSongs    int     `json:"info total songs"`
	TotalArtists  int     `json:"info total artists"`
	TotalAlbums   int     `json:"info total albums"`
	TotalGenres   int     `json:"info total genres"`
	TotalDuration float64 `json:"info total duration"`
}

type PlayerDetails struct {
	Power       int    `json:"power"`
	IsPlaying   int    `json:"isplaying"`
	DisplayType string `json:"displaytype"`
	Model       string `json:"model"`
	ModelName   string `json:"modelname"`
	UUID        string `json:"uuid"`
	Connected   int    `json:"connected"`
	IsPlayer    int    `json:"isplayer"`
	IP          string `json:"ip"`
	CanPowerOff int    `json:"canpoweroff"`
	Name        string `json:"name"`
	PlayerID    string `json:"playerid"`
	Firmware    any `json:"firmware"`
}

type ExtendedServerStatus struct {
	ServerStatus
	Players []PlayerDetails `json:"players_loop"`
}

type PlayerStatus struct {
	PlayerID string // not set by JSON, set by client method
		
	SyncMaster string `json:"sync_master"`
	SyncSlaves string `json:"sync_slaves"`
	
	Power       int    `json:"power"`
	PlayerName string `json:"player_name"`
	Mode string `json:"mode"`
	PlayerConnected int `json:"player_connected"`
	SignalStrength int `json:"signalstrength"`
	
	PlaylistCurIndex any `json:"playlist_cur_index"`
	PlaylistRepeat int `json:"playlist_repeat"`
	PlaylistTimestamp float64 `json:"playlist_timestamp"`
	PlaylistShuffle int `json:"playlist_shuffle"`
	PlaylistMode string `json:"playlist_mode"`
	PlaylistTracks int `json:"playlist_tracks"`
	IP string `json:"player_ip"`

	MixerVolume int `json:"mixer volume"`
}