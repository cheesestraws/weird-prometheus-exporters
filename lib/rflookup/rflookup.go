package rflookup

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

/* This is a horrible hack to look up friendly names of RF sensors
   per RTL_433.  In DNS, because I'm a masochist.  Don't use this
   unless you know what you're doing.  Vaguely, at least.

   We look up a CNAME for
      <model>-<label>-<value>-<label>...<value>.<frequency-in-mhz>.rf

   The list of labels to be included is the following.  Add new labels
   to the end of the list only, to avoid breaking existing hostnames:
*/

var rflookupLabels []string = []string{
	"channel",
	"id",
	"data",
}

/* This is the error to proagate upwards if we get no model for some
   reason */

var ErrNoModel error = errors.New("no model provided for sensor")

/* If the CNAME does not end in .friendly.rf, we return this error */

var ErrWrongDomain error = errors.New("returned CNAME not in .friendly.rf")

func MkLookupHostname(labelset map[string]string, centreFrequency int) (string, error) {
	// Do we have a model?
	_, ok := labelset["model"]
	if !ok {
		return "", ErrNoModel
	}

	freqString := fmt.Sprintf("%d", centreFrequency/1000000)

	// construct hosrtname
	fqdn := labelset["model"]
	for _, l := range rflookupLabels {
		if val, ok := labelset[l]; ok {
			fqdn += fmt.Sprintf("-%s-%s", l, val)
		}
	}

	fqdn += fmt.Sprintf(".%s.rf", freqString)
	fqdn = strings.ToLower(fqdn)

	return fqdn, nil
}

func Lookup(labelset map[string]string, centreFrequency int) (string, error) {
	fqdn, err := MkLookupHostname(labelset, centreFrequency)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	hn, err := net.DefaultResolver.LookupCNAME(ctx, fqdn)
	if err != nil {
		return hn, err
	}
	
	if !strings.HasSuffix(hn, ".friendly.rf.") {
		return hn, ErrWrongDomain
	}
	
	return strings.TrimSuffix(hn, ".friendly.rf."), nil
}
