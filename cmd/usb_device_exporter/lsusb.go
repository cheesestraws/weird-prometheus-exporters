package main

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
var usbLineRE = regexp.MustCompile(`^Bus ([0-9a-fA-F]+) Device ([0-9a-fA-F]+): ID ([0-9a-fA-F]+):([0-9a-fA-F]+) (.*)$`)

type USBDevice struct {
	Bus         string `prometheus_label:"bus"`
	Device      string `prometheus_label:"device"`
	ID          string `prometheus_label:"id"`
	Description string `prometheus_label:"description"`
}

var ErrBadLine = errors.New("bad line")

func parseUSBDevice(s string) (USBDevice, error) {
	fields := usbLineRE.FindStringSubmatch(s)
	if fields == nil {
		return USBDevice{}, ErrBadLine
	}

	return USBDevice{
		Bus:         fields[1],
		Device:      fields[2],
		ID:          fields[3] + ":" + fields[4],
		Description: fields[5],
	}, nil
}

type USBDevices struct {
	Devices           map[USBDevice]int `prometheus_map:"device"`
	DeviceFetchErrors int               `prometheus:"device_fetch_errors"`
}

func parseUSBDevices(s string) USBDevices {
	devs := USBDevices{
		Devices:           make(map[USBDevice]int),
		DeviceFetchErrors: 0,
	}

	ss := strings.Split(s, "\n")
	for _, line := range ss {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		dev, err := parseUSBDevice(strings.TrimSpace(line))
		if err != nil {
			devs.DeviceFetchErrors++
		} else {
			devs.Devices[dev]++
		}
	}

	return devs
}

func lsusb() USBDevices {
	bs, err := exec.Command("lsusb").Output()
	if err != nil {
		log.Printf("err: %v", err)
		return USBDevices{
			DeviceFetchErrors: 1,
		}
	}

	s := string(bs)

	return parseUSBDevices(s)
}
