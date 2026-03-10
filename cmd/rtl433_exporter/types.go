package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"
)

type JSONNumberOrString struct {
	IsNumber bool
    Number float64
    IsString bool
    String string
    
    IsOther bool
}

func (j *JSONNumberOrString) UnmarshalJSON(b []byte) error {
	// try number first
	err := json.Unmarshal(b, &j.Number)
	if err == nil {
		j.IsString = false
		j.IsNumber = true
		return nil
	}
	
	err = json.Unmarshal(b, &j.String)
	if err == nil {
		j.IsString = true
		j.IsNumber = false
		return nil
	}
	
	j.IsOther = true
	
	return nil
}

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

type LabelSet string
type DynamicSensorMetrics struct {
	LastSeen time.Time
	Metrics map[string]float64
}

type AllDynamicSensorMetrics struct {
	Metrics map[LabelSet]DynamicSensorMetrics
}

func (a *AllDynamicSensorMetrics) Observe(labels map[string]string, metrics map[string]float64) {
	if a.Metrics == nil {	
		// lazily initialise map if we need to
		a.Metrics = make(map[LabelSet]DynamicSensorMetrics)
	}
	
	var ll []string
	for k, v := range labels {
		ll = append(ll, fmt.Sprintf("%s=%q", k, v))
	}
	slices.Sort(ll)
	
	ls := strings.Join(ll, ",")
	
	dsm := DynamicSensorMetrics{
		LastSeen: time.Now(),
		Metrics: maps.Clone(metrics),
	}
	
	a.Metrics[LabelSet(ls)] = dsm
}

func (a *AllDynamicSensorMetrics) PromBytes(prefix string, baseURL string) []byte {
	var accum bytes.Buffer
	
	for labels, dsm := range a.Metrics {
		fmt.Fprintf(&accum, prefix + "timestamp{base_url=\"%s\",%s} %d\n", baseURL, labels, dsm.LastSeen.Unix())
		for k, v := range dsm.Metrics {
			fmt.Fprintf(&accum, prefix + "%s{base_url=\"%s\",%s} %v\n", k, baseURL, labels, v)
		}
	}
	
	return accum.Bytes()
}

func (a *AllDynamicSensorMetrics) FlushOldCrap(timeout time.Duration) {
	if a.Metrics == nil {
		return
	}
	
	for labels, dsm := range a.Metrics {
		if time.Since(dsm.LastSeen) > timeout {
			delete(a.Metrics, labels)
		}
	}
}

type Metrics struct {
	sync.Mutex

	Metadata              map[PromMetadata]int `prometheus_map:"metadata"`
	MetadataValid         int                  `prometheus:"metadata_valid"`
	MetadataPollSuccesses int                  `prometheus:"metadata_poll_success_count"`
	MetadataPollFailures  int                  `prometheus:"metadata_poll_failure_count"`
	
	StreamConnectionUp int `prometheus:"stream_connection_up"`
	
	DynamicMetrics AllDynamicSensorMetrics
	LastDynamicMetricFlush int64 `prometheus:"last_dynamic_metric_flush"`
}
