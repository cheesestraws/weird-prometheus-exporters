package macutils

import (
	"runtime"
	"testing"
)

func TestVersionRoughly(t *testing.T) {
	ver := OSVersion()
	if runtime.GOOS != "darwin" && ver != nil {
		t.Errorf("got version %v for OS %v", ver, runtime.GOOS)
	}
	if runtime.GOOS == "darwin" && ver == nil {
		t.Errorf("got nil version for OS %v", runtime.GOOS)
	}
	if ver != nil {
		t.Logf("you're running on %v", ver)
	}
}
