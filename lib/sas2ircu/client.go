package sas2ircu

import (
	"strings"
	"regexp"
)

type SAS2IRCUAdapter struct {
	Index string
	AdapterType string
	PCIAddress string
}

var listLine = regexp.MustCompile(`^(\d+)\s+(\w+)\s+\w+\s+\w+\s+([0-9a-fA-Fh:]+)`)

func parseSAS2IRCUList(output []byte) map[string]SAS2IRCUAdapter {
	ret := make(map[string]SAS2IRCUAdapter)
	lines := strings.Split(string(output), "\n")
	
	for _, l := range lines {
		l = strings.TrimSpace(l)
		
		matches := listLine.FindStringSubmatch(l)
		if matches == nil {
			continue
		}
		
		stats := SAS2IRCUAdapter{
			Index: matches[1],
			AdapterType: matches[2],
			PCIAddress: matches[3],
		}
		ret[matches[1]] = stats
	}
	
	return ret
}

var extractThreeLetterCode = regexp.MustCompile(`\(([A-Z]+)\)`)

type IRStatus int
const (
	IRStatusUnknown IRStatus = iota
	IRStatusOkay
	IRStatusDegraded
	IRStatusFailed
	IRStatusMissing
	IRStatusInitializing
	IRStatusOnline
)

func irStatusFromString(s string) IRStatus {
	matches := extractThreeLetterCode.FindStringSubmatch(s)
	if matches == nil {
		return IRStatusUnknown
	}
	
	switch matches[1] {
		case "OKY": return IRStatusOkay
		case "DGD": return IRStatusDegraded
		case "FLD": return IRStatusFailed
		case "MIS": return IRStatusMissing
		case "INIT": return IRStatusInitializing
		case "ONL": return IRStatusOnline
	}
	return IRStatusUnknown
}

type IR struct {
	VolumeID string
	Status IRStatus
}

type PDStatus int
const (
	PDStatusUnknown PDStatus = iota
	PDStatusOnline
	PDStatusHotSpare
	PDStatusReady
	PDStatusAvailable
	PDStatusFailed
	PDStatusMissing
	PDStatusStandby
	PDStatusOutOfSync
	PDStatusDegraded
	PDStatusRebuilding
	PDStatusOptimal
)

func pdStatusFromString(s string) PDStatus {
	matches := extractThreeLetterCode.FindStringSubmatch(s)
	if matches == nil {
		return PDStatusUnknown
	}
	
	switch matches[1] {
		case "ONL": return PDStatusOnline
		case "HSP": return PDStatusHotSpare
		case "RDY": return PDStatusReady
		case "AVL": return PDStatusAvailable
		case "FLD": return PDStatusFailed
		case "MIS": return PDStatusMissing
		case "SBY": return PDStatusStandby
		case "OSY": return PDStatusOutOfSync
		case "DGD": return PDStatusDegraded
		case "RBLD": return PDStatusRebuilding
		case "OPT": return PDStatusOptimal

	}
	return PDStatusUnknown
}


type PhysicalDevice struct {
	Enclosure string
	Slot string
	State string
	SerialNumber string
	Protocol string
	DriveType string
}