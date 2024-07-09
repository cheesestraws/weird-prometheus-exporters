package main

import (
	"fmt"
	"strings"

	"github.com/cheesestraws/weird-prometheus-exporters/lib/fn"
	"github.com/cheesestraws/weird-prometheus-exporters/lib/netatalk"
)

type NBPEntity struct {
	Object string
	Type string
	Zone string
}

func StringToEntity(s string) NBPEntity {
	ot, z, _ := strings.Cut(s, "@")
	o, t, _ := strings.Cut(ot, ":")
	
	return NBPEntity{o, t, z}
}

type NetworkState struct {
	QueryError bool
	
	Zones []string
	ZoneAppleRouters map[string]int
	ZoneEntityCount map[string]int
	
	NBPEntities []NBPEntity
}

func ParseNetworkState(zones []string, entities map[string]map[string]string) *NetworkState {
	ns := NetworkState{
		QueryError: false,
		Zones: zones,
		ZoneAppleRouters: make(map[string]int),
		ZoneEntityCount: make(map[string]int),
	}
	
	for _, zone := range zones {
		// Get the entities
		entities := fn.Mapkeymap(entities[zone], func(s string) NBPEntity {
			e := StringToEntity(s)
			e.Zone = zone
			return e
		})
		
		// Count 'em
		ns.ZoneEntityCount[zone] = len(entities)
		
		// AppleRouters?
		routers := fn.Filter(entities, func(e NBPEntity) bool { return e.Type == "AppleRouter" })
		ns.ZoneAppleRouters[zone] = len(routers)
		
		ns.NBPEntities = append(ns.NBPEntities, entities...)
	}
	
	return &ns
}

func QueryNetworkState() (*NetworkState, error) {
	nserr := &NetworkState{QueryError: true}
	
	zs, err := netatalk.DefaultQuerier.GetZones()
	if err != nil {
		return nserr, err
	}
	
	zentities := make(map[string]map[string]string)
	
	for _, zone := range zs {
		zentities[zone], err = netatalk.DefaultQuerier.NBPLookup(fmt.Sprintf("@%s", zone))
		if err != nil {
			return nserr, err
		}
	}
	
	return ParseNetworkState(zs, zentities), nil
}

