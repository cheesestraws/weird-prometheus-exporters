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
	DeviceIsA string
	Enclosure string
	Slot string
	State PDStatus
	SerialNumber string
	Protocol string
	DriveType string
}

type Devices struct {
	IRs []IR
	PhysicalDevices []PhysicalDevice
}

var dataLineRegex = regexp.MustCompile(`^\s+([^\s].*[^\s])\s+:\s+(.*)$`)
var irLineRegex = regexp.MustCompile(`^IR volume \d`)
var pdLineRegex = regexp.MustCompile(`^Device is a (.*)$`)
var dashedLineRegex = regexp.MustCompile(`^----`)

func parseSAS2IRCUDisplay(output []byte) Devices {
	var physDevices []PhysicalDevice
	var irs []IR
	
	kvs := make(map[string]string)
	var isIRVolume bool
	var isPhysicalDevice bool
	
	flush := func() {		
		if isPhysicalDevice {
			pd := PhysicalDevice{
				DeviceIsA: kvs["DeviceIsA"],
				Enclosure: kvs["Enclosure #"],
				Slot: kvs["Slot #"],
				State: pdStatusFromString(kvs["State"]),
				SerialNumber: kvs["Serial No"],
				Protocol: kvs["Protocol"],
				DriveType: kvs["Drive Type"],
			}
			physDevices = append(physDevices, pd)
		}
		
		if isIRVolume {
			ir := IR{
				VolumeID: kvs["Volume ID"],
				Status: irStatusFromString(kvs["Status of volume"]),
			}
			irs = append(irs, ir)
		}
		
			
		isIRVolume = false
		isPhysicalDevice = false
		kvs = make(map[string]string)
	}
	
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
	
		// is it an IR header line?
		matches := irLineRegex.FindStringSubmatch(line)
		if matches != nil {
			flush()
			isIRVolume = true
			continue
		}
		
		matches = pdLineRegex.FindStringSubmatch(line)
		if matches != nil {
			flush()
			isPhysicalDevice = true
			kvs["DeviceIsA"] = matches[1]
		}
		
		// is it a data line?
		matches = dataLineRegex.FindStringSubmatch(line)
		if matches != nil {
			kvs[matches[1]] = matches[2]
		}
	}
	
	flush()
	
	return Devices{irs, physDevices}
}