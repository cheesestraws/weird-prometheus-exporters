package main

import (
	"fmt"
	"log"
	"net/url"
	"time"
)

func fetchOneMetadata(url string) (RPCDeviceInfo, RPCMeta, error) {
	var devs RPCDeviceInfo
	err := doJSONRPCQuery(
		url,
		"get_dev_info",
		&devs,
	)
	if err != nil {
		return RPCDeviceInfo{}, RPCMeta{}, err
	}

	var meta RPCMeta
	err = doJSONRPCQuery(
		url,
		"get_meta",
		&meta,
	)
	if err != nil {
		return RPCDeviceInfo{}, RPCMeta{}, err
	}

	return devs, meta, nil
}

func fetchMetadata(baseURL string) {
	url, err := url.JoinPath(baseURL, "/jsonrpc")
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(10 * time.Second)
		devs, meta, err := fetchOneMetadata(url)
		metrics.Lock()
		if err != nil {
			log.Printf("metadata poll err: %v", err)

			metrics.MetadataValid = 0
			metrics.MetadataPollFailures++
		} else {
			metrics.MetadataValid = 1
			metrics.MetadataPollSuccesses++
			metrics.Metadata = map[PromMetadata]int{
				metadataRPCToProm(devs, meta): 1,
			}
		}
		metrics.Unlock()

		if *dump {
			fmt.Printf("%s\n\n\n", getBody())
		}
	}
}
