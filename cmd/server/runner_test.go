package main

import (
	"testing"
)

func TestConvertResult(t *testing.T) {
	input := "1|label1:value1,label2:value2"

	value, labels := convertResult(input)

	if value != 1 {
		t.Errorf("Expected metric value %f but found %f", 1.0, value)
	}
	if len(labels) != 2 {
		t.Errorf("Expected %d labels but found %d", 2, len(labels))
	}
	if labels["label2"] != "value2" {
		t.Errorf("Expected label to contain %s but found %s", "value2", labels["label2"])
	}
}
