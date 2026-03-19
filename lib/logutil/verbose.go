package logutil

import (
	"flag"
	"fmt"
	"log"
)

var verbose *bool

func Verbose(s string) {
	if *verbose {
		log.Print(s)
	}
}

func Verbosef(format string, args ...any) {
	Verbose(fmt.Sprintf(format, args...))
}

func init() {
	verbose = flag.Bool("v", false, "verbose logging")
}
