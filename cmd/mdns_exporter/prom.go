package main

type Service struct {
	Name string `prometheus_label:"name"`
	Desc string `prometheus_label:"description"`
}

type Output struct {
	ServiceLastSeen map[Service]int `prometheus_map:"service_last_seen"`
}
