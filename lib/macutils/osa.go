package macutils

import (
	"strings"
	"os/exec"
)

func EscapeAppleScriptString(s string) string {
	// Need to replace " and \, and that's it
	return strings.ReplaceAll(
		strings.ReplaceAll(
			s, `\`, `\\`,
		),
		`"`, `\"`,
	)	
}

func ExecuteAppleScript(script ...string) ([]byte, error) {
	var args []string
	
	for _, v := range script {
		args = append(args, "-e", v)
	}
	
	output, err := exec.Command("osascript", args...).Output()
	return output, err
}