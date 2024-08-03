package main

import (
	"net"
)

type ResponseType string

type logEntry struct {
	IP       net.IP
	Type     RRType
	Response ResponseType
	Hostname string
}
