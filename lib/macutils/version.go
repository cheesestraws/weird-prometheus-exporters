package macutils

import (
	"os/exec"
	"runtime"
	"strings"
)

func OSVersion() []string {
	if runtime.GOOS != "darwin" {
		return nil
	}

	output, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return nil
	}

	return strings.Split(strings.TrimSpace(string(output)), ".")
}
