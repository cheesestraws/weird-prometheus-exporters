package macutils

import (
	"runtime"
	"testing"
)

func maconly(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skipf("this test is mac only")
	}
}
