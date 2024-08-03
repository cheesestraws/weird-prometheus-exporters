package main

import (
	"log"
	"strings"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

func firstfield(line string) fn.Maybe[string] {
	if len(line) <= 1 {
		return fn.Absent[string]()
	}

	fld, _, _ := strings.Cut(line[1:], ":")
	return fn.Present(fld)
}

type DatabaseSummary struct {
}

func checkData(dataFile string) (interface{}, error) {
	return nil, nil
}

func getSOAs(lines []string) map[string]struct{} {
	domains := make(map[string]struct{})

	for _, line := range lines {
		// is this an SOA type line?
		if strings.HasPrefix(line, ".") {
			domain := firstfield(line)
			domain.Range(func(domain string) {
				// if we have a domain, add
				domains[domain] = struct{}{}
			})
		}
	}

	log.Printf("%+v", domains)
	return domains
}

func checkSOAUsefulness(lines []string) {
}
