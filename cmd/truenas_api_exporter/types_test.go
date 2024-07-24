package main

import "testing"

func TestBytes(t *testing.T) {
	b, ok := stringToBytes("1G / fish funky stuff")
	if !ok {
		t.Errorf("no return")
	}
	if b != 1024*1024*1024 {
		t.Errorf("where's my gig?  got %v", b)
	}

	b, ok = stringToBytes("15T / fish funky stuff")
	if !ok {
		t.Errorf("no return")
	}
	if b != 1024*1024*1024*1024*15 {
		t.Errorf("where's my gig?  got %v", b)
	}

}
