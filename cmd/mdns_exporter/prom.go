package main

type Service struct {
	Name string `prometheus_label:"service"`
	Desc string `prometheus_label:"description"`
}

type Resource struct {
	Service string `prometheus_label:"service"`
	Name    string `prometheus_label:"name"`
	Host    string `prometheus_label:"host"`
	AddrV4  string `prometheus_label:"addr_v4"`
	AddrV6  string `prometheus_label:"addr_v6"`
	Port    int    `prometheus_label:"port"`
	Info    string `prometheus_label:"info"`
}

type Metrics struct {
	RunningServicePollers int64 `prometheus:"running_service_pollers"`
	ServicePollCount      int64 `prometheus:"service_poll_count"`
	ServiceReplyCount     int64 `prometheus:"service_reply_count"`
	ServiceCleanupCount   int64 `prometheus:"service_cleanup_count"`

	RunningResourcePollers int64 `prometheus:"running_resource_pollers"`
	ResourcePollCount      int64 `prometheus:"resource_poll_count"`
	ResourceReplyCount     int64 `prometheus:"resource_reply_count"`
	ResourceCleanupCount   int64 `prometheus:"resource_cleanup_count"`

	ServiceLastSeen  map[Service]int64  `prometheus_map:"service_last_seen"`
	ResourceLastSeen map[Resource]int64 `prometheus_map:"resource_last_seen"`
}
