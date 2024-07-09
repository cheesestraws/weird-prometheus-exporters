package netatalk

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

// A Querier can go and get information about AppleTalk zones and NBP
// entries.
type Querier interface {
	GetZones() ([]string, error)
	NBPLookup(pattern string) (map[string]string, error)
}

type cliQuerier struct{}

var DefaultQuerier Querier = cliQuerier{}

func (c cliQuerier) GetZones() ([]string, error) {
	output, err := exec.Command("getzones").Output()
	if err != nil {
		return nil, err
	}
	outstr := string(bytes.TrimSpace(output))

	return strings.Split(outstr, "\n"), nil
}

var spaces *regexp.Regexp = regexp.MustCompile("\\s+")

func (c cliQuerier) NBPLookup(pattern string) (map[string]string, error) {
	output, err := exec.Command("getzones").Output()
	if err != nil {
		return nil, err
	}
	outstr := string(bytes.TrimSpace(output))
	lines := strings.Split(outstr, "\n")

	m := make(map[string]string)

	for _, l := range lines {
		ll := strings.TrimSpace(l)
		flds := spaces.Split(ll, -1)

		if len(flds) != 2 {
			return nil, errors.New("invalid formatting from nmblkup")
		}

		m[flds[0]] = flds[1]
	}

	return m, nil
}
