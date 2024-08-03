package main

import (
	"net"
)

type RRType int

var rrTypes = map[int]string{
	0x01: "A",
	0x02: "NS",
	0x05: "CNAME",
	0x06: "SOA",
	0x0c: "PTR",
	0x0f: "MX",
	0x10: "TXT",
	0x1c: "AAAA",
	0x26: "A6",
	0xfb: "IXFR",
	0xfc: "AXFR",
	0xff: "wildcard",
}

type ResponseType string

type logEntry struct {
	IP       net.IP
	Type     RRType
	Response ResponseType
	Hostname string
}
