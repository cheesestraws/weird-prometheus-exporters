package main

import (
	"fmt"
	"sync"
)

type RPCDeviceInfo struct {
	Vendor  string
	Product string
	Serial  string
}

type RPCMeta struct {
	Frequencies        []int
	HopTimes           []int `json:"hop_times"`
	CentreFrequency    int   `json:"center_frequency"`
	Duration           int
	SampleRate         int `json:"samp_rate"`
	ConversionMode     int `json:"conversion_mode"`
	FSKPulseDetectMode int `json:"fsk_pulse_detect_mode"`
}

type PromMetadata struct {
	Vendor             string `prometheus_label:"vendor"`
	Product            string `prometheus_label:"product"`
	Serial             string `prometheus_label:"serial"`
	Frequencies        string `prometheus_label:"frequencies"`
	HopTimes           string `prometheus_label:"hop_times"`
	CentreFrequency    int    `prometheus_label:"centre_frequency"`
	Duration           int    `prometheus_label:"duration"`
	SampleRate         int    `prometheus_label:"sample_rate"`
	ConversionMode     int    `prometheus_label:"conversion_mode"`
	FSKPulseDetectMode int    `prometheus_label:"fsk_pulse_detect_mode"`
}

func metadataRPCToProm(devs RPCDeviceInfo, meta RPCMeta) PromMetadata {
	return PromMetadata{
		Vendor:             devs.Vendor,
		Product:            devs.Product,
		Serial:             devs.Serial,
		Frequencies:        fmt.Sprint(meta.Frequencies),
		HopTimes:           fmt.Sprint(meta.HopTimes),
		CentreFrequency:    meta.CentreFrequency,
		Duration:           meta.Duration,
		SampleRate:         meta.SampleRate,
		ConversionMode:     meta.ConversionMode,
		FSKPulseDetectMode: meta.FSKPulseDetectMode,
	}
}

type Metrics struct {
	sync.Mutex

	Metadata              map[PromMetadata]int `prometheus_map:"metadata"`
	MetadataValid         int                  `prometheus:"metadata_valid"`
	MetadataPollSuccesses int                  `prometheus:"metadata_poll_success_count"`
	MetadataPollFailures  int                  `prometheus:"metadata_poll_failure_count"`
}
