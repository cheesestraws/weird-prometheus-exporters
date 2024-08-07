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