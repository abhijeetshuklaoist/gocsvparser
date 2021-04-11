package utils

import (
	"csvpaser/parser"
	"csvpaser/utils"
	"strings"
	"testing"
)

func TestConvertToLowerCase(t *testing.T) {
	testString := "MixedCaseString"
	total :=  utils.ConvertToLowerCaseString(testString)
	if total != strings.ToLower(testString) {
		t.Errorf("Sum was incorrect, got: %s, want: %s.", total, "mixedcasestring")
	}
}

func TestConvertToString(t *testing.T) {
	testString := "\uFEFFMixedCaseString"
	total := utils.ConvertToLowerCaseString(testString)
	if total != strings.ToLower(strings.ReplaceAll(testString, "\uFEFF","")) {
		t.Errorf("Sum was incorrect, got: %s, want: %s.", total, "mixedcasestring")
	}
}

func TestFetchHeaderIndex(t *testing.T) {
	testHeaders := make(map[string]int)
	testHeaders["TestHeader1"] = 0
	testHeaders["TestHeader2"] = 1
	testHeaders["TestHeader3"] = 2

	actualHeaderIndex, err := parser.FetchHeaderIndex(testHeaders, "TestHeader3")
	expectedHeaderIndex := 2
	if err != nil {
		t.Errorf("Header %s was expected to be found in %v", "TestHeader3", testHeaders)
	}
	if actualHeaderIndex != expectedHeaderIndex {
		t.Errorf("Header %s was expected to be found in %v at index %d", "TestHeader3", testHeaders, expectedHeaderIndex)
	}
}

func TestFetchHeaderIndexForNonExistingHeader(t *testing.T) {
	testHeaders := make(map[string]int)
	testHeaders["TestHeader1"] = 0

	actualHeader, err := parser.FetchHeaderIndex(testHeaders, "TestHeader2")
	if err == nil {
		t.Errorf("Header %s was not expected to be found in %v", "TestHeader2", testHeaders)
	}
	if actualHeader != -1 {
		t.Errorf("Header %s was expected to be found in %v at index %d", "TestHeader3", testHeaders, -1)
	}
}