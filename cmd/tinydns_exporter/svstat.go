package main

import (
	"os/exec"
	"regexp"
	"strconv"
)

type ServiceStatus struct {
	Disappeared int `prometheus:"svc_disappeared"`
	Up          int `prometheus:"svc_up"`
	Time        int `prometheus:"svc_state_time"`
	PID         int `prometheus:"svc_pid"`
}

var statusRE = regexp.MustCompile(`^[^:]*:\s*(\S+)`)
var pidRE = regexp.MustCompile(`\(pid ([0-9]+)`)
var uptimeRE = regexp.MustCompile(`([0-9]+) seconds`)

func svstatResult(result string) ServiceStatus {
	var ss ServiceStatus

	statusSlice := statusRE.FindStringSubmatch(result)
	if statusSlice == nil {
		ss.Disappeared = 1
	}
	if statusSlice != nil {
		if statusSlice[1] == "up" {
			ss.Up = 1
		} else {
			ss.Up = 0
		}
	}

	pidSlice := pidRE.FindStringSubmatch(result)
	if pidSlice != nil {
		ss.PID, _ = strconv.Atoi(pidSlice[1])
	}

	uptimeSlice := uptimeRE.FindStringSubmatch(result)
	if uptimeSlice != nil {
		ss.Time, _ = strconv.Atoi(uptimeSlice[1])
	}

	return ss
}

func Svstat(path string) (ServiceStatus, error) {
	out, err := exec.Command("svstat", path).Output()
	if err != nil {
		return ServiceStatus{}, err
	}

	return svstatResult(string(out)), nil
}
