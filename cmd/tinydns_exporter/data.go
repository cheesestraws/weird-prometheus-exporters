package main

import (
	"context"
	"net"
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
	DomainError: "Error",
	DomainNoNSes: "No NS records",
	DomainHasOurNS: "Our NS records",
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
