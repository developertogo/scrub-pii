package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
)

func TestScrubPII(t *testing.T) {
	RegisterTestingT(t)
	var input interface{}
	var expected interface{}
	var actual interface{}

	// unit tests all sub-folders under the tests dir
	dir := "./tests/"
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		fmt.Println("Test ", f.Name())
		t.Log("*********************************************************************************")
		t.Log("Test ", f.Name())
		t.Log("*********************************************************************************")

		testDir := dir + f.Name()
		inputPath := fmt.Sprintf("%s/input.json", testDir)
		outputPath := fmt.Sprintf("%s/output.json", testDir)
		sensitiveFieldsPath := fmt.Sprintf("%s/sensitive_fields.txt", testDir)

		inputByte, _ := ioutil.ReadFile(inputPath)
		_ = json.Unmarshal([]byte(inputByte), &input)
		expectedByte, _ := ioutil.ReadFile(outputPath)
		_ = json.Unmarshal([]byte(expectedByte), &expected)

		// scrub it
		actualStr, _ := ScrubPersonalInfo(inputPath, sensitiveFieldsPath)
		_ = json.Unmarshal([]byte(actualStr), &actual)

		t.Log("input:    ", input)
		t.Log("expected: ", expected)
		t.Log("output:   ", actual)

		Expect(expected).To(Equal(actual))
		fmt.Println("... Pass")
	}
}
