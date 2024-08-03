package main

import (
	"context"
	"os"
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
		"example.com":  struct{}{},
		"example2.com": struct{}{},
	}) {
		t.Errorf("getSOAs returned bad data")
	}
}

func TestGetNS(t *testing.T) {
	nses, err := getNS(context.Background(), "example.com")
	if err != nil {
		t.Errorf("getNS shouldn't return error %v", err)
	}
	if len(nses) == 0 {
		t.Errorf("example.com should always have nameservers")
	}

	nses, _ = getNS(context.Background(), "example.invalid")
	if len(nses) != 0 {
		t.Errorf("anything under .invalid should never have nameservers")
	}
}

func TestDomainStatus(t *testing.T) {
	s := domainStatus(context.Background(), "example.invalid", ".example.com")
	if s != DomainError {
		t.Errorf(".invalid domains shouldn't have NSes")
	}

	s = domainStatus(context.Background(), "example.com", ".iana-servers.net")
	if s != DomainHasOurNS {
		t.Errorf("example.com should have IANA nameservers")
	}

	s = domainStatus(context.Background(), "example.com", ".ns.ecliptiq.co.uk")
	if s != DomainDoesNotHaveOurNS {
		t.Errorf("example.com should not have ecliptiq nameservers")
	}
}

func TestSOANSes(t *testing.T) {
	testFile := strings.Split(`
.example.com::a.ns.ecliptiq.co.uk:360
&example.com::b.ns.ecliptiq.co.uk:360
+foo.example.com:192.168.0.1
+bar.example.com:192.168.0.2
+woo.example.com:192.168.0.3
.foo.invalid::a.ns.ecliptiq.co.uk:360
&foo.invalid::b.ns.ecliptiq.co.uk:360
`, "\n")

	ds := checkSOANS(context.Background(), testFile, ".iana-servers.net")
	expected := map[string]DomainStatus{
		"example.com": DomainHasOurNS,
		"foo.invalid": DomainError,
	}

	if !reflect.DeepEqual(ds, expected) {
		t.Errorf("checkSOANS didn't return expected value")
	}
}

func TestCheckData(t *testing.T) {
	dataFile := os.Getenv("DATAFILE")
	suffix := os.Getenv("SUFFIX")

	if dataFile == "" || suffix == "" {
		t.Skipf("to manually check CheckData, set DATAFILE and SUFFIX")
	}

	d, err := checkData(context.Background(), dataFile, suffix)
	if err != nil {
		t.Errorf("checkData failed: %v", err)
	}

	t.Logf("%+v", d)
}
