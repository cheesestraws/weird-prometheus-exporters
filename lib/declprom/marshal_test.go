package declprom

import (
	"testing"
)

func TestLabelsRoughlyWork(t *testing.T) {
	m := Marshaller{
		BaseLabels: map[string]string{
			"a": "b",
			"c": "d",
		},
	}

	s := m.metricLabelString(nil, map[string]string{"e": "f"}, map[string]string{"g": "h"})
	if s != `a="b",c="d",e="f",g="h"` {
		t.Errorf("labels are broken")
	}
}

func TestMarshalRoughlyWorks(t *testing.T) {
	s := struct {
		A int     `prometheus:"E"`
		B float64 `prometheus:"F"`

		C map[string]int `prometheus_map:"G" prometheus_map_key:"K"`

		unexported int `prometheus:"ignored"`
	}{
		A: 1,
		B: 2.2,
		C: map[string]int{
			"hello": 3,
			"world": 4,
		},
	}

	m := Marshaller{
		MetricNamePrefix: "test_",
		BaseLabels: map[string]string{
			"one": "2",
			"two": "1",
		},
	}
	t.Logf("%s", m.Marshal(s, nil))
}
