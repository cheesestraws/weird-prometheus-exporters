package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFirstField(t *testing.T) {
	line := "+"
	fld := firstfield(line)
	fld.Range(func(s string) {
		t.Errorf("Expected to get empty field")
	})

	line = ".foo.com"
	fld = firstfield(line)
	fld.OrElse(func() {
		t.Errorf("should not have got empty field")
	})
	fld.Range(func(s string) {
		if s != "foo.com" {
			t.Errorf("wrong field contents: %s", s)
		}
	})

	line = "+host.example.com:123"
	fld = firstfield(line)
	fld.OrElse(func() {
		t.Errorf("should not have got empty field")
	})
	fld.Range(func(s string) {
		if s != "host.example.com" {
			t.Errorf("wrong field contents: %s", s)
		}
	})
}

func TestGetSOAs(t *testing.T) {
	testFile := strings.Split(`
.example.com::a.ns.ecliptiq.co.uk:360
&example.com::b.ns.ecliptiq.co.uk:360
+foo.example.com:192.168.0.1
+bar.example.com:192.168.0.2
+woo.example.com:192.168.0.3
.example2.com::a.ns.ecliptiq.co.uk:360
&example2.com::b.ns.ecliptiq.co.uk:360
`, "\n")
	soas := getSOAs(testFile)
	if !reflect.DeepEqual(soas, map[string]struct{}{
		"example.com": struct{}{},
		"example2.com": struct{}{},
	}) {
		t.Errorf("getSOAs returned bad data")
	}

}
