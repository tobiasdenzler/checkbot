package main

import (
	"fmt"
	"testing"
)

func TestFindLineInFile(t *testing.T) {
	line, err := findLineInFile("../../scripts/compliance/missing_quota_on_project_total.sh", "# TYPE")

	if err != nil {
		t.Error(err)
	}

	if line != "" {
		fmt.Println(line)
	} else {
		t.Error("empty result")
	}
}
