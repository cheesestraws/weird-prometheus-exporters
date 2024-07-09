package main

import (
	"fmt"
)

func main() {
	ns, err := QueryNetworkState()
	fmt.Printf("err: %v\n", err)
	fmt.Printf("ns: %+v\n", *ns)
	fmt.Printf("prom: %s\n", ns.ToPrometheus())
}