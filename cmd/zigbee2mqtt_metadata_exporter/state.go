package main

import (
	"sync"
)

type Device struct {
	FriendlyName   string `prometheus_label:"sensor"`
	IEEEAddress    string `prometheus_label:"ieee_address"`
	Type           string `prometheus_label:"type"`
	ModelID        string `prometheus_label:"model_id"`
	NetworkAddress string `prometheus_label:"network_address"`
	Description    string `prometheus_label:"description"`
	Model          string `prometheus_label:"model"`
	Vendor         string `prometheus_label:"vendor"`
}

type Link struct {
	Relationship string `prometheus_label:"relationship"`
	SourceSensor string `prometheus_label:"source_sensor"`
	SourceIEEEAddr string `prometheus_label:"source_ieee_address"`
	SourceType string `prometheus_label:"source_type"`
	TargetSensor string `prometheus_label:"target_sensor"`
	TargetIEEEAddr string `prometheus_label:"target_ieee_address"`
	TargetType string `prometheus_label:"target_type"`
}

type State struct {
	sync.Mutex
	
	DevicesTimestamp int64 `prometheus:"device_info_timestamp"`
	Devices map[Device]int `prometheus_map:"device_info"`
	
	LinksTimestamp int64 `prometheus:"link_info_timestamp"`
	Links map[Link]int `prometheus_map:"link_info"`
	
	ExporterPublishErrors int `prometheus:"exporter_publish_errors"`
	ExporterPublishSuccess int `prometheus:"exporter_publish_success"`
}

var promstate State
