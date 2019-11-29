package main

import (
	"reflect"
	"testing"
	"time"
)

type testpairResult struct {
	input  string
	value  float64
	labels map[string]string
}

var testsResult = []testpairResult{
	{"1|label1=value1", 1, map[string]string{"label1": "value1"}},
	{"1|label1=value1,label2=value2", 1, map[string]string{"label1": "value1", "label2": "value2"}},
	{"1|label1=value1,label2=value2,label3=value3", 1, map[string]string{"label1": "value1", "label2": "value2", "label3": "value3"}},
	{"1", 1, make(map[string]string)},
	{"1|user=system:admin", 1, map[string]string{"user": "system:admin"}},
}

func TestConvertResult(t *testing.T) {
	for _, pair := range testsResult {
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

type testpairFile struct {
	path     string
	filename string
	result   string
	hasError bool
}

var testFile = []testpairFile{
	{"../../test/scripts/single_result.sh", "single_result", "42|label1=value1,label2=value2\n", false},
	{"../../test/scripts/failed_result.sh", "failed_result", "", true},
	{"../../test/scripts/empty_result.sh", "empty_result", "\n", false},
}

func TestRunScript(t *testing.T) {

	for _, pair := range testFile {

		file := pair.path
		check := new(Check)
		check = &Check{
			Name:        pair.filename,
			File:        file,
			Interval:    10,
			Active:      true,
			MetricType:  extractMetadataFromFile(metaType, file),
			Help:        extractMetadataFromFile(metaHelp, file),
			stoppedchan: make(chan struct{}),
			Nextrun:     time.Now().Unix(),
		}

		result, err := runBashScript(*check)

		if err != nil {
			if !pair.hasError {
				t.Error("Error happened: ", err)
			}
		}
		if result != pair.result {
			t.Errorf("Expected result is %s but got %s", pair.result, result)
		}
	}
}

func TestRegisterMetricsSuccess(t *testing.T) {

	check := getPlaceholderCheck("test_valid", "Gauge")

	// Register a valid metric
	registerMetricsForCheck(check, 42, map[string]string{"label1": "value1"})
}

func TestRegisterMetricsFailed(t *testing.T) {

	check := getPlaceholderCheck("test_invalid", "Gauge")

	// You can only register metrics with the same labels, this is not valid!
	registerMetricsForCheck(check, 42, map[string]string{"label1": "value1"})
	registerMetricsForCheck(check, 43, map[string]string{"label2": "value2"})
}

func getPlaceholderCheck(metricName string, metricType string) *Check {

	check := new(Check)
	check = &Check{
		Name:        metricName,
		File:        "placeholder",
		Interval:    10,
		Active:      true,
		MetricType:  metricType,
		Help:        "placeholder",
		stoppedchan: make(chan struct{}),
		Nextrun:     time.Now().Unix(),
	}

	return check
}
