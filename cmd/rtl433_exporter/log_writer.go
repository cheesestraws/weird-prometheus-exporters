package main

import (
	"bytes"
	"os"
)

// a Writer that filters out annoying error messages
type logWriter struct {}

func (l logWriter) Write(buf []byte) (int, error) {
	// "Unsolicited response received on idle HTTP channel" is due to some
	// weird behaviour in rtl_433; nowt we can do about it but suppress
	// the message.
	if bytes.Contains(buf, []byte("Unsolicited response received on idle HTTP channel")) {
		return len(buf), nil
	}

	return os.Stderr.Write(buf)
}
