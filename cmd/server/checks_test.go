package main

import (
	"fmt"
	"testing"
)

func TestFindLineInFile(t *testing.T) {
	line, err := findLineInFile("../../test/scripts/single_result.sh", "# TYPE")

	if err != nil {
		t.Error(err)
	}

	if line != "" {
		fmt.Println(line)
	} else {
		t.Error("empty result")
	}
}
