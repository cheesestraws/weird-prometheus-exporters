package main

import (
	"strconv"
)

type Device struct {
	Name string
	DeviceDiscoveryID string
	RawDeviceModel string
	BatteryStatus string
	BatteryLevel float64
	DeviceStatus string
	ThisDevice bool
	IsMac bool
}


// A Dev is a key for a map of device data
type Dev struct {
	Name string `prometheus_label:"device_name"`
	RawDeviceModel string `prometheus_label:"device_model"`
	ThisDevice bool `prometheus_label:"exporting_device"`
	IsMac bool `prometheus_label:"is_mac"`
	DeviceDiscoveryID string `prometheus_label:"device_discovery_id"`
}

type DeviceMetrics struct {
	Online map[Dev]int `prometheus_map:"online"`
	BatteryStatus map[Dev]int `prometheus_map:"battery_status" prometheus_help:"1 => charging, 0 => not charging, -1 => unknown"`
	BatteryLevel map[Dev]float64 `prometheus_map:"battery_level"`
	DeviceStatus map[Dev]int `prometheus_map:"device_status"`
}

func deviceMetricsFromDevices(ds []Device) DeviceMetrics {
	d := DeviceMetrics{
		Online: make(map[Dev]int),
		BatteryStatus: make(map[Dev]int),
		BatteryLevel: make(map[Dev]float64),
		DeviceStatus: make(map[Dev]int),
	}
	
	for _, v := range ds {
		dev := Dev{
			Name: v.Name,
			RawDeviceModel: v.RawDeviceModel,
			ThisDevice: v.ThisDevice,
			IsMac: v.IsMac,
			DeviceDiscoveryID: v.DeviceDiscoveryID,
		}
		
		// Online is a wild guess
		if v.DeviceStatus == "200" {
			d.Online[dev] = 1
		} else {
			d.Online[dev] = 0
		}
		
		if v.BatteryStatus == "Charging" {
			d.BatteryStatus[dev] = 1
		} else if v.BatteryStatus == "NotCharging" {
			d.BatteryStatus[dev] = 0
		} else {
			d.BatteryStatus[dev] = -1
		}
		
		d.BatteryLevel[dev] = v.BatteryLevel
		
		status, err := strconv.Atoi(v.DeviceStatus)
		if err != nil {
			status = -1
		}
		d.DeviceStatus[dev] = status
	}
	
	return d
}