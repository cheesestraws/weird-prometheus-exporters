package main

import "testing"

func TestTinyDNSLineToRRTypes(t *testing.T) {
	rrs := TinyDNSLineToRRTypes(":speakersassociates.com:16:\214v=spf1\040include\072eu.zcsend.net\040include\072_spf.google.com\040a\072mail.speakersassociates.com\040a\072c-mx0.midworld.co.uk\040include\072spf.mta01.mailhawk.io\040-all:3600")
	if rrs[0] != RRTypeTXT {
		t.Errorf("Got %v, expected TXT", rrs[0])
	}
}
