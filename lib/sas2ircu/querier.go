package sas2ircu

import (
	"os/exec"
)

type Querier struct {
	Executable string
}

type Response struct {
	Adapters map[string]Adapter
	Devices map[string]Devices
}

func (q *Querier) getBytes(args ...string) ([]byte, error) {
	return exec.Command(q.Executable, args...).Output()
}

func (q *Querier) Get() (*Response, error) {
	r := &Response{}
	
	bs, err := q.getBytes("LIST")
	if err != nil {
		return nil, err
	}
	
	r.Adapters = parseSAS2IRCUList(bs)
	r.Devices = make(map[string]Devices)
	
	for id := range r.Adapters {
		bs, err := q.getBytes(id, "DISPLAY")
		if err != nil {
			return nil, err
		}
		
		devs := parseSAS2IRCUDisplay(bs)
		r.Devices[id] = devs
	}
	
	return r, nil
}