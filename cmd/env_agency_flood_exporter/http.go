package main

import (
	"net/http"
	"encoding/json"
	"io"
	"fmt"
	"time"
)

func makeHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func fetch(cli *http.Client, stationID int) error {
	url := fmt.Sprintf("https://environment.data.gov.uk/flood-monitoring/id/stations/%d", stationID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var sta station
	err = json.Unmarshal(body, &sta)
	if err != nil {
		return err
	}
	
	fmt.Printf("%+v", riverLevelFromStation(&sta))
	return nil
}
