package macutils

import (
	"testing"
)

func TestEscaping(t *testing.T) {
	s := EscapeAppleScriptString(`"hello\"`)
	if s != `\"hello\\\"` {
		t.Errorf("applescript escaping is broken")
	}
}

func TestExecuteAppleScript(t *testing.T) {
	maconly(t)

	out, err := ExecuteAppleScript(
		`tell application "Finder"`,
		`end tell`,
		`copy "the wombles of wimbledon" to stdout`,
	)
	
	if err != nil {
		t.Errorf("ExecuteAppleScript returned error %v", err)
	}
	
	o := string(out)
	if o != "the wombles of wimbledon\n" {
		t.Errorf("got unexpected output %+v", o)
	}
}