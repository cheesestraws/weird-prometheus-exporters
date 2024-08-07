package main

import (
	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
	"github.com/cheesestraws/weird-prometheus-exporters/lib/sas2ircu"
)

type Adapter struct {
	Index       string `prometheus_label:"adapter_index"`
	AdapterType string `prometheus_label:"adapter_type"`
	PCIAddress  string `prometheus_label:"pci_address"`
}

type IR struct {
	AdapterIndex string `prometheus_label:"adapter_index"`
	VolumeID     string `prometheus_label:"volume_id"`
}

type PhysicalDevice struct {
	AdapterIndex string `prometheus_label:"adapter_index"`
	DeviceIsA    string `prometheus_label:"device_is_a"`
	Enclosure    string `prometheus_label:"enclosure"`
	Slot         string `prometheus_label:"slot"`
	SerialNumber string `prometheus_label:"device_serial_number"`
	Protocol     string `prometheus_label:"protocol"`
	DriveType    string `prometheus_label:"drive_type"`
}

type Summary struct {
	AdapterDetails       map[Adapter]int                      `prometheus_map:"adapter_details"`
	IRStatus             map[IR]sas2ircu.IRStatus             `prometheus_map:"ir_status"`
	PhysicalDeviceStatus map[PhysicalDevice]sas2ircu.PDStatus `prometheus_map:"physical_device_status"`

	FetchError int `prometheus:"in_fetch_error_state"`
}

func Summarise(r *sas2ircu.Response) *Summary {
	sums := &Summary{
		IRStatus:             make(map[IR]sas2ircu.IRStatus),
		PhysicalDeviceStatus: make(map[PhysicalDevice]sas2ircu.PDStatus),
	}

	sums.AdapterDetails = fn.Mapmap(r.Adapters,
		func(_ string, a sas2ircu.Adapter) (Adapter, int) {
			return Adapter{
				Index:       a.Index,
				AdapterType: a.AdapterType,
				PCIAddress:  a.PCIAddress,
			}, 1
		})

	for adapter_index, devs := range r.Devices {
		// IRs first
		for _, ir := range devs.IRs {
			myir := IR{
				AdapterIndex: adapter_index,
				VolumeID:     ir.VolumeID,
			}
			sums.IRStatus[myir] = ir.Status
		}

		// then physdevs
		for _, pd := range devs.PhysicalDevices {
			mydev := PhysicalDevice{
				AdapterIndex: adapter_index,
				DeviceIsA:    pd.DeviceIsA,
				Enclosure:    pd.Enclosure,
				Slot:         pd.Slot,
				SerialNumber: pd.SerialNumber,
				Protocol:     pd.Protocol,
				DriveType:    pd.DriveType,
			}
			sums.PhysicalDeviceStatus[mydev] = pd.State
		}
	}

	return sums
}
