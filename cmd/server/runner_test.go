package main

import (
	"reflect"
	"testing"
)

type testpair struct {
	input  string
	value  float64
	labels map[string]string
}

var tests = []testpair{
	{"1|label1:value1,label2:value2", 1, map[string]string{"label1": "value1", "label2": "value2"}},
	{"1|label1:value1", 1, map[string]string{"label1": "value1"}},
	{"1", 1, make(map[string]string)},
	{"1|user:system:admin", 1, map[string]string{"user": "system:admin"}},
}

func TestConvertResult(t *testing.T) {
	for _, pair := range tests {
		value, labels := convertResult(pair.input)

		if value != pair.value {
			t.Errorf("Expected metric value %f but found %f", pair.value, value)
		}
		if len(labels) != len(pair.labels) {
			t.Errorf("Expected %d labels but found %d", len(pair.labels), len(labels))
		}
		eq := reflect.DeepEqual(labels, pair.labels)
		if !eq {
			t.Errorf("Expected labels %s but found %s", labels, pair.labels)
		}
	}
}
