package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/declprom"
)

func Test_parseUSBDevice(t *testing.T) {
	// bad line
	_, err := parseUSBDevice("farts")
	if err == nil {
		t.Fatalf("parsed a fart")
	}

	// ok line

	dev, err := parseUSBDevice("Bus 001 Device 002: ID 2109:3431 VIA Labs, Inc. Hub")
	if err != nil {
		t.Fatalf("couldn't parse line")
	}

	expectedDev := USBDevice{
		Bus:         "001",
		Device:      "002",
		ID:          "2109:3431",
		Description: "VIA Labs, Inc. Hub",
	}

	if !reflect.DeepEqual(dev, expectedDev) {
		t.Fatalf("got bad answer")
	}
}

func Test_parseUSBDevices(t *testing.T) {
	cmdOut := `Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
honks
Bus 001 Device 002: ID 2109:3431 VIA Labs, Inc. Hub
Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
farts
`

	devs := parseUSBDevices(cmdOut)

	expectedDevs := USBDevices{
		Devices: map[USBDevice]int{
			USBDevice{Bus: "001", Device: "001", ID: "1d6b:0002", Description: "Linux Foundation 2.0 root hub"}: 1,
			USBDevice{Bus: "001", Device: "002", ID: "2109:3431", Description: "VIA Labs, Inc. Hub"}:            1,
			USBDevice{Bus: "002", Device: "001", ID: "1d6b:0003", Description: "Linux Foundation 3.0 root hub"}: 1,
		},
		DeviceFetchErrors: 2,
	}

	if !reflect.DeepEqual(devs, expectedDevs) {
		t.Fatalf("got bad answer")
	}

	// Does it serialise?
	m := declprom.Marshaller{
		MetricNamePrefix: "usb_",
	}

	bs := m.Marshal(devs, nil)
	s := string(bs)

	expectedProm := `# TYPE usb_device gauge
usb_device{bus="002",description="Linux Foundation 3.0 root hub",device="001",id="1d6b:0003"} 1
usb_device{bus="001",description="VIA Labs, Inc. Hub",device="002",id="2109:3431"} 1
usb_device{bus="001",description="Linux Foundation 2.0 root hub",device="001",id="1d6b:0002"} 1
# TYPE usb_device_fetch_errors gauge
usb_device_fetch_errors{} 2`

	if !reflect.DeepEqual(strings.TrimSpace(s), expectedProm) {
		t.Fatalf("bad serialisation")
	}
}
