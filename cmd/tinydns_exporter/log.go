package main

import (
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/nxadm/tail"
)

var loglineRE = regexp.MustCompile(`([0-9a-fA-F]+):[0-9a-fA-F]+:[0-9a-fA-F]+ (.) ([0-9a-fA-F]+) ([^ ]+)$`)

func responseTypeToLongerString(r string) string {
	switch r {
	case "+":
		return "ok"
	case "-":
		return "dropped"
	case "I":
		return "not impl"
	case "C":
		return "bad class"
	case "/":
		return "garbled"
	default:
		return "???"
	}
}

type DomainTypeResponse struct {
	Domain   string `prometheus_label:"domain"`
	Type     RRType `prometheus_label:"type"`
	Response string `prometheus_label:"response"`
}

type LogStats struct {
	RequestCount map[DomainTypeResponse]int `prometheus_map:"request_count"`
}

func NewLogStats() LogStats {
	return LogStats{
		RequestCount: make(map[DomainTypeResponse]int),
	}
}

func updateStats(hexIP string, responseType string, rrType RRType, hostname string) {
	// do we have SOAs?
	var domains map[string]struct{}
	state.l.RLock()
	state.soas.Range(func(m map[string]struct{}) {
		domains = m
	})
	state.l.RUnlock()

	// find the relevant domain
	var domain string
	if domains != nil {
		domain = findSuffix(hostname, domains).Or("(unknown)")
	} else {
		domain = "(domain-info-unavailable)"
	}

	key := DomainTypeResponse{
		Domain:   domain,
		Type:     rrType,
		Response: responseType,
	}

	state.l.Lock()
	state.logStats.RequestCount[key]++
	state.l.Unlock()

	if *verbose {
		log.Printf("%s, %s, %v, %s, %s", hexIP, responseType, rrType, hostname, domain)
	}
}

func watchLogs() error {
	t, err := tail.TailFile(
		filepath.Join(*logdir, "current"),
		tail.Config{
			Follow:        true,
			ReOpen:        true,
			CompleteLines: true,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: io.SeekEnd,
			},
		},
	)

	if err != nil {
		return err
	}

	for line := range t.Lines {
		parts := loglineRE.FindStringSubmatch(line.Text)
		if parts != nil {
			hexIP := parts[1]

			responseType := responseTypeToLongerString(parts[2])
			rrTypeInt, _ := strconv.ParseInt(parts[3], 16, 32)
			rrType := RRType(rrTypeInt)
			hostname := parts[4]

			updateStats(hexIP, responseType, rrType, hostname)
		}
	}

	return t.Err()
}
