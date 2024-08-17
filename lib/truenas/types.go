package truenas

type Alert struct {
	UUID string `json:"uuid"`
	Source string `json:"source"`
	Dismissed bool `json:"dismissed"`
	Level string `json:"level"`
	FomattedText string `json:"formatted"`
}

var AlertValidLevels []string = []string{"INFO", "NOTICE", "WARNING", 
	"ERROR", "CRITICAL", "ALERT", "EMERGENCY"}
	
type Device struct {
	Type string `json:"type"`
	
	Stats struct {
		ReadErrors int `json:"read_errors"`
		WriteErrors int `json:"write_errors"`
		ChecksumErrors int `json:"checksum_errors"`
		Size int `json:"size"`
		Allocated int `json:"allocated"`
	}`json:"stats"`
}

type Pool struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Status string `json:"status"`
	
	Topology struct {
		Data []Device `json:"data"`
	} `json:"topology"`
	
	Healthy bool `json:"healthy"`
	
}

var PoolValidStatuses []string = []string{"ONLINE", "DEGRADED", "FAULTED",
	"OFFLINE", "UNAVAIL", "REMOVED"}
	
type CloudSync struct {
	ID int `json:"id"`
	Description string `json:"description"`
	Enabled bool `json:"enabled"`
	Path string `json:"path"`
	Job struct {
		State string `json:"state"`
		Progress struct {
			Percent int `json:"percent"`
			Description string `json:"description"`
		} `json:"progress"`
	} `json:"job"`
}

var CloudSyncStatuses []string = []string{"FAILED", "ABORTED",
	"PENDING", "RUNNING", "SUCCESS"}