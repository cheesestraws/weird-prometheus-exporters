package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
)

type DomainStatus int

const (
	DomainUnknown DomainStatus = iota
	DomainError
	DomainNoNSes
	DomainHasOurNS
	DomainDoesNotHaveOurNS
	DomainStatusEnd
)

var dsStrings = map[DomainStatus]string{
	DomainError:            "Error",
	DomainNoNSes:           "No NS records",
	DomainHasOurNS:         "Our NS records",
	DomainDoesNotHaveOurNS: "Only other NS records",
}

func (d DomainStatus) String() string {
	s, ok := dsStrings[d]
	if !ok {
		return "Unknown"
	}
	return s
}

func firstfield(line string) fn.Maybe[string] {
	if len(line) <= 1 {
		return fn.Absent[string]()
	}

	fld, _, _ := strings.Cut(line[1:], ":")
	return fn.Present(fld)
}

type DomainAndStatus struct {
	Domain string       `prometheus_label:"domain"`
	Status DomainStatus `prometheus_label:"status"`
}

type DomainAndType struct {
	Domain string `prometheus_label:"domain"`
	Type   RRType `prometheus_label:"type"`
}

type DatabaseSummary struct {
	DomainStatus  map[DomainAndStatus]int `prometheus_map:"domain_status"`
	StatusSummary map[DomainStatus]int    `prometheus_map:"status_summary" prometheus_map_key:"status"`
	RecordTypes   map[DomainAndType]int   `prometheus_map:"record_types"`
}

func getSOAsFromFile(dataFile string) (map[string]struct{}, error) {
	bs, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bs), "\n")

	return getSOAs(lines), nil
}

func checkData(ctx context.Context, dataFile string, ourSuffix string) (DatabaseSummary, error) {
	var s DatabaseSummary

	bs, err := os.ReadFile(dataFile)
	if err != nil {
		return s, err
	}

	lines := strings.Split(string(bs), "\n")

	// Check SOA statuses
	statuses := checkSOANS(ctx, lines, ourSuffix)
	s.DomainStatus = fn.Mapmap(statuses,
		func(k string, v DomainStatus) (DomainAndStatus, int) {
			return DomainAndStatus{k, v}, 1
		})
	s.StatusSummary = fn.CountMapValues(statuses)
	s.RecordTypes = summariseRecordTypes(lines)

	return s, nil
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
	return domains
}

func getNS(ctx context.Context, domain string) ([]string, error) {
	r := &net.Resolver{}

	nses, err := r.LookupNS(ctx, domain)
	if err != nil {
		return nil, err
	}

	return fn.Map(nses, func(ns *net.NS) string { return ns.Host }), nil
}

func domainStatus(ctx context.Context, domain string, ourSuffix string) DomainStatus {
	if *verbose {
		log.Printf("checking NS records for %v", domain)
	}
	nses, err := getNS(ctx, domain)
	if err != nil {
		return DomainError
	}
	if len(nses) == 0 {
		return DomainNoNSes
	}

	for _, ns := range nses {
		if strings.HasSuffix(strings.TrimSuffix(ns, "."), ourSuffix) {
			return DomainHasOurNS
		}
	}

	return DomainDoesNotHaveOurNS
}

func checkSOANS(ctx context.Context, lines []string, ourSuffix string) map[string]DomainStatus {
	soas := getSOAs(lines)
	return fn.Mapmap(soas, func(k string, _ struct{}) (string, DomainStatus) {
		return k, domainStatus(ctx, k, ourSuffix)
	})
}

func findSuffix(s string, suffixes map[string]struct{}) fn.Maybe[string] {
	for suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return fn.Present(suffix)
		}
	}

	return fn.Absent[string]()
}

func summariseRecordTypes(lines []string) map[DomainAndType]int {
	m := make(map[DomainAndType]int)
	domains := getSOAs(lines)

	for _, line := range lines {
		// Do we care about this line?
		rrs := TinyDNSLineToRRTypes(line)
		if len(rrs) == 0 {
			continue
		}

		// Do we even have an FQDN?
		firstfield(line).Range(func(fqdn string) {
			// What domain does this line belong to?

			domain := findSuffix(fqdn, domains).Or("(unattached)")

			// increment counters for each rr type
			for _, rr := range rrs {
				m[DomainAndType{domain, rr}]++
			}
		})
	}

	return m
}
