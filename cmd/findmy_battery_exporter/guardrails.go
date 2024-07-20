package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/macutils"
)

func needAccessTo(fn string) {
	// attempt to read the file
	_, err := os.ReadFile(fn)
	if err != nil {
		log.Printf("couldn't open necessary file: do I have full disk access?")
		os.Exit(1)
	}
}

func checkForAccess() {
	home := os.Getenv("HOME")
	if home == "" {
		log.Printf("$HOME not set")
		os.Exit(1)
	}
	
	cachedir := filepath.Join(home, "Library/Caches/com.apple.findmy.fmipcore")
	
	needAccessTo(filepath.Join(cachedir, "Devices.data"))
}

func openFindMyApp() {
	_, err := macutils.ExecuteAppleScript(`tell application "FindMy" to launch`)
	if err != nil {
		log.Printf("couldn't open Find My: %v", err)
		os.Exit(1)
	}
}

func warnForUntestedVersions() {
	v := macutils.OSVersion()
	
	if v == nil {
		log.Printf("couldn't determine macOS version")
		os.Exit(1)
	}
	
	if v[0] != "12" && v[0] != "13" {
		log.Printf("warning: you're running on an untested macOS version.")
		log.Printf("warning: this will certainly not work > 14.3.1, and may break in other places.")
	}
}