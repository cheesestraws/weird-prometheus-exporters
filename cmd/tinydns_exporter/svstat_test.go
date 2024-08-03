package main

import (
	"reflect"
	"testing"
)

func TestSvstatResult(t *testing.T) {
	var ss ServiceStatus
	var expected ServiceStatus

	expected = ServiceStatus{Up: 1, Time: 13234356, PID: 29170}
	ss = svstatResult("/service/tinydns: up (pid 29170) 13234356 seconds")
	if !reflect.DeepEqual(ss, expected) {
		t.Errorf("up didn't work")
	}

	expected = ServiceStatus{Up: 0, Time: 1234, PID: 0}
	ss = svstatResult("/service/tinydns: down 1234 seconds")
	if !reflect.DeepEqual(ss, expected) {
		t.Errorf("down didn't work")
	}

}
