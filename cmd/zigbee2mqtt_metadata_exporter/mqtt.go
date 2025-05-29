package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func relationshipString(relationship int) string {
	switch(relationship) {
	case 0x0:
		return "parent"
	case 0x1:
		return "child"
	case 0x2:
		return "sibling"
	case 0x3:
		return "none_of_the_above"
	case 0x4:
		return "previous_child"
	default:
		return "unknown"
	}
}

func handleMapMessage(client mqtt.Client, message mqtt.Message) {
	type jsonNode struct {
		FriendlyName   string `json:"friendlyName"`
		IEEEAddress    string `json:"ieeeAddr"`
		Type string `json:"type"`
	}
	
	type jsonLink struct {
		Relationship int `json:"relationship"`
		SourceIEEEAddr string `json:"sourceIeeeAddr"`
		TargetIEEEAddr string `json:"targetIeeeAddr"`
	}
	
	type jsonResponse struct {
		Data struct {
			Value struct {
				Nodes []jsonNode `json:"nodes"`
				Links []jsonLink `json:"links"`
			} `json:"value"`
		} `json:"data"`
	}
	
	var resp jsonResponse
	err := json.Unmarshal(message.Payload(), &resp)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
	
	nodeMap := make(map[string]jsonNode)
	for _, n := range resp.Data.Value.Nodes {
		nodeMap[n.IEEEAddress] = n
	}
	
	promLinks := make(map[Link]int)
	for _, l := range resp.Data.Value.Links {
		promLinks[Link{
			Relationship: relationshipString(l.Relationship),
			SourceSensor: nodeMap[l.SourceIEEEAddr].FriendlyName,
			SourceIEEEAddr: l.SourceIEEEAddr,
			SourceType:nodeMap[l.SourceIEEEAddr].Type,
			TargetSensor: nodeMap[l.TargetIEEEAddr].FriendlyName,
			TargetIEEEAddr: l.TargetIEEEAddr,
			TargetType: nodeMap[l.TargetIEEEAddr].Type,
		}] = 1
	}
	
	promstate.Lock()
	defer promstate.Unlock()
	promstate.Links = promLinks
	promstate.LinksTimestamp = time.Now().UnixMilli()
	
	log.Printf("got %d links", len(promLinks))
}

func handleDevicesMessage(client mqtt.Client, message mqtt.Message) {
	type jsonDevice struct {
		FriendlyName   string `json:"friendly_name"`
		IEEEAddress    string `json:"ieee_address"`
		Type           string `json:"type"`
		ModelID        string `json:"model_id"`
		NetworkAddress int    `json:"network_address"`
		Description    string `json:"description"`
		Definition     struct {
			Model  string `json:"model"`
			Vendor string `json:"vendor"`
		} `json:"definition"`
	}

	var ds []jsonDevice

	err := json.Unmarshal(message.Payload(), &ds)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}

	promDevices := make(map[Device]int)
	log.Printf("got %d devices", len(ds))
	for _, d := range ds {
		promDevices[Device{
			FriendlyName:   d.FriendlyName,
			IEEEAddress:    d.IEEEAddress,
			Type:           d.Type,
			ModelID:        d.ModelID,
			NetworkAddress: fmt.Sprintf("%v", d.NetworkAddress),
			Description:    d.Description,
			Model:          d.Definition.Model,
			Vendor:         d.Definition.Vendor,
		}] = 1
	}
	
	promstate.Lock()
	defer promstate.Unlock()
	promstate.Devices = promDevices
	promstate.DevicesTimestamp = time.Now().UnixMilli()
}

func mapRequestLoop(client mqtt.Client, base string, everyN int) {
	mapReq := path.Join(base, "bridge/request/networkmap")
	payload := `{"type": "raw", "routes": false}`
	for {
		token := client.Publish(mapReq, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			log.Printf("[err] %v", token.Error())
			promstate.Lock()
			promstate.ExporterPublishErrors++
			promstate.Unlock()
		} else {
			promstate.Lock()
			promstate.ExporterPublishSuccess++
			promstate.Unlock()
		}
		
		if (everyN > 0) {
			time.Sleep(time.Duration(everyN) * time.Hour)
		} else {
			select{}
		}
	}
}

func run_mqtt(broker string, user string, pass string, base string, everyN int) {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	hostname, _ := os.Hostname()
	opts := mqtt.NewClientOptions()
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.SetClientID(hostname + strconv.Itoa(time.Now().Second()))
	opts.SetOrderMatters(false)
	opts.AddBroker(broker)

	cli := mqtt.NewClient(opts)
	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		log.Printf("Connected to %s\n", broker)
	}

	// Subscribe to topics
	devicesTopic := path.Join(base, "bridge/devices")
	if token := cli.Subscribe(devicesTopic, 0, handleDevicesMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	
	mapTopic := path.Join(base, "bridge/response/networkmap")
	if token := cli.Subscribe(mapTopic, 0, handleMapMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mapRequestLoop(cli, base, everyN)
}
