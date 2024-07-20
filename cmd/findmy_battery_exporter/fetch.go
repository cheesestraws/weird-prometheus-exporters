package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func fetchDevices() ([]Device, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return nil, errors.New("HOME unset")
	}

	cachedir := filepath.Join(home, "Library/Caches/com.apple.findmy.fmipcore")
	j, err := os.ReadFile(filepath.Join(cachedir, "Devices.data"))
	if err != nil {
		return nil, err
	}

	var ds []Device
	err = json.Unmarshal(j, &ds)

	return ds, err
}
