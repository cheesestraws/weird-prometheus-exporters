package main

import (
	"fmt"
	"strings"
	"bytes"

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

func (n NBPEntity) AsLabels() string {
	return fmt.Sprintf("object=%q,type=%q,zone=%q",
		n.Object, n.Type, n.Zone)
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

func (ns *NetworkState) ToPrometheus() []byte {
	bs := &bytes.Buffer{}
	
	help := func(metric string, help string) {
		fmt.Fprintf(bs, "# HELP %s %s\n", metric, help)
	}
	
	mtype := func(metric string, t string) {
		fmt.Fprintf(bs, "# TYPE %s %s\n", metric, t)
	}
	
	val := func(metric string, i int) {
		fmt.Fprintf(bs, "%s %d\n", metric, i)
	}

	// QueryError
	help("appletalk_query_error", "NBP failed while looking up network state")
	mtype("appletalk_query_error", "gauge")
	
	var qe int 
	if ns.QueryError {
		qe = 1
	}
	
	val("appletalk_query_error", qe)
	
	// Zones
	help("appletalk_zone", "AppleTalk zones seen")
	mtype("appletalk_zone", "gauge")
	for _, z := range ns.Zones {
		m := fmt.Sprintf("appletalk_zone{zone=%q}", z)
		val(m, 1)		
	}
	
	// Apple Routers
	help("appletalk_zone_applerouters", "number of Apple routers seen")
	mtype("appletalk_zone_applerouters", "gauge")
	for _, z := range ns.Zones {
		m := fmt.Sprintf("appletalk_zone_applerouters{zone=%q}", z)
		val(m, ns.ZoneAppleRouters[z])
	}
	
	// Entity count
	help("appletalk_zone_entity_count", "number of entities in zone")
	mtype("appletalk_zone_entity_count", "gauge")
	for _, z := range ns.Zones {
		m := fmt.Sprintf("appletalk_zone_entity_count{zone=%q}", z)
		val(m, ns.ZoneEntityCount[z])
	}

	// Entities
	help("appletalk_nbp_entity", "nbp entities")
	mtype("appletalk_nbp_entity", "gauge")
	for _, e := range ns.NBPEntities {
		m := fmt.Sprintf("appletalk_nbp_entity{%s}", e.AsLabels())
		val(m, 1)
	}
	
	return bs.Bytes()
}

