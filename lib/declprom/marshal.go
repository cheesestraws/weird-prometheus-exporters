package declprom

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

type RenamePrometheusMetricer interface {
	RenamePrometheusMetric(baseName string) string
}

type PrometheusLabelser interface {
	PrometheusLabels() map[string]string
}

type Marshaller struct {
	MetricNamePrefix string
	MetricNameSuffix string

	BaseLabels map[string]string
}

func (m *Marshaller) metricLabels(i interface{}, localLabels map[string]string, extraLabels map[string]string) map[string]string {
	labelMap := make(map[string]string)

	// Can we get labels from the struct?
	lr, ok := i.(PrometheusLabelser)
	if ok {
		sl := lr.PrometheusLabels()
		for k, v := range sl {
			labelMap[k] = v
		}
	}

	// Get labels from the marshaller
	for k, v := range m.BaseLabels {
		labelMap[k] = v
	}

	// Get labels from caller
	for k, v := range localLabels {
		labelMap[k] = v
	}

	// Get labels from fields (like maps)
	for k, v := range extraLabels {
		labelMap[k] = v
	}

	return labelMap
}

func (m *Marshaller) metricLabelString(i interface{}, localLabels map[string]string, extraLabels map[string]string) string {
	var labels []string
	lmap := m.metricLabels(i, localLabels, extraLabels)

	for k, v := range lmap {
		labels = append(labels, fmt.Sprintf("%s=%q", k, v))
	}

	sort.Strings(labels)

	return strings.Join(labels, ",")
}

func (m *Marshaller) metricName(baseName string, i interface{}) string {
	name := baseName

	// Can we rewrite?
	re, ok := i.(RenamePrometheusMetricer)
	if ok {
		name = re.RenamePrometheusMetric(baseName)
	}

	return m.MetricNamePrefix + name + m.MetricNameSuffix
}

func (m *Marshaller) taggedStructToLabels(s reflect.Value) map[string]string {
	accum := make(map[string]string)
	myType := s.Type()
	for i := 0; i < myType.NumField(); i++ {
		fld := myType.Field(i)
		if !fld.IsExported() {
			continue
		}
		
		lname := fld.Tag.Get("prometheus_label")
		val := fmt.Sprintf("%v", s.Field(i).Interface())
		accum[lname]=val
	}
	return accum
}

func (m *Marshaller) Marshal(s interface{}, localLabels map[string]string) []byte {
	bs := &bytes.Buffer{}

	mtype := func(metric string, t string) {
		fmt.Fprintf(bs, "# TYPE %s %s\n", metric, t)
	}
	
	mhelp := func(metric string, h string) {
		if h != "" {
			fmt.Fprintf(bs, "# HELP %s %s\n", metric, h)
		}
	}


	val := func(metric string, labels string, val interface{}) {
		var valstr string
		switch v := val.(type) {
		case time.Duration:
			valstr = fmt.Sprintf("%v", v.Seconds())
		default:
			valstr = fmt.Sprintf("%v", v)
		}
		fmt.Fprintf(bs, "%s{%s} %s\n", metric, labels, valstr)
	}

	// Now iterate through the fields
	theStruct := reflect.ValueOf(s)
	myType := theStruct.Type()

	for i := 0; i < myType.NumField(); i++ {
		fld := myType.Field(i)
		if !fld.IsExported() {
			continue
		}
		
		help := fld.Tag.Get("prometheus_help")

		// is it a scalar?
		basename := fld.Tag.Get("prometheus")
		if basename != "" {
			realName := m.metricName(basename, theStruct)
			labels := m.metricLabelString(
				theStruct,
				localLabels,
				nil,
			)
			value := theStruct.Field(i).Interface()
			mtype(realName, "gauge")
			mhelp(realName, help)
			val(realName, labels, value)
		}

		// is it a map?
		mapname := fld.Tag.Get("prometheus_map")
		keyname := fld.Tag.Get("prometheus_map_key")
		// map with a struct key?
		if mapname != "" && fld.Type.Kind() == reflect.Map && fld.Type.Key().Kind() == reflect.Struct {
			// Iterate over it
			realName := m.metricName(mapname, theStruct)

			mtype(realName, "gauge")
			mhelp(realName, help)
			iter := theStruct.Field(i).MapRange()
			for iter.Next() {
				keyLabels := m.taggedStructToLabels(iter.Key())
				labels := m.metricLabelString(
					theStruct,
					localLabels,
					keyLabels,
				)
				metricValue := iter.Value().Interface()
				val(realName, labels, metricValue)
			}
		}
		if mapname != "" && keyname != "" && fld.Type.Kind() == reflect.Map {
			// Iterate over it
			realName := m.metricName(mapname, theStruct)

			mtype(realName, "gauge")
			mhelp(realName, help)
			iter := theStruct.Field(i).MapRange()
			for iter.Next() {
				tagValue := fmt.Sprintf("%v", iter.Key().Interface())
				labels := m.metricLabelString(
					theStruct,
					localLabels,
					map[string]string{
						keyname: tagValue,
					},
				)
				metricValue := iter.Value().Interface()
				val(realName, labels, metricValue)
			}
		}
	}

	return bs.Bytes()
}
