package utils

import (
	"csvpaser/utils"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateCSV(t *testing.T) {
	testCSV := createTempDir() + "/" + "temp_test.csv"
	_, _ = utils.CreateCSV(testCSV)

	path, _ := filepath.Abs(testCSV)
	_, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if errors.Is(err, os.ErrNotExist) {
		t.Errorf("Could not create csv %v", path)
	}
	removeTempDir()
}


func TestCWriteDataInCSV(t *testing.T) {
	testCSV := createTempDir() + "/" + "temp_test.csv"
	_, _ = utils.CreateCSV(testCSV)

	dataChannel := make(chan []string, 1)
	boolChannel := make(chan bool, 1)
	expectedData := [] string {"Col1", "Col2"}
	dataChannel <- expectedData
	go utils.WriteDataInCSV(testCSV, dataChannel, boolChannel)

	actualData := utils.ReadDataFromCSV(testCSV)

	for _, data := range actualData {
		if !Equal(data, expectedData) {
			t.Errorf("Wrong data persisted/read from csv, expecting: %s, actual: %s.", expectedData, data)
		}
	}
	close(dataChannel)
	<- dataChannel
	<- boolChannel
	close(boolChannel)
	removeTempDir()
}

func TestCWriteDataInCSVSync(t *testing.T) {
	testCSV := createTempDir() + "/" + "temp_test2.csv"
	_, _ = utils.CreateCSV(testCSV)
	expectedData := [] string {"Col1","Col2"}
	utils.WriteDataInCSVSync(testCSV, expectedData)
	actualData := utils.ReadDataFromCSV(testCSV)
	for _, data := range actualData {
		if !Equal(data, expectedData) {
			t.Errorf("Wrong data persisted/read from csv, expecting: %s, actual: %s.", expectedData, data)
		}
	}
	removeTempDir()
}

func createTempDir() string {
	parentDir := os.TempDir()
	tempDir, err := ioutil.TempDir(parentDir, "csv_util_test")
	if err != nil {
		log.Fatal(err)
	}
	return tempDir
}

func removeTempDir() {
	parentDir := os.TempDir()
	globPattern := filepath.Join(parentDir, "csv_util_test*")
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		log.Fatalf("Failed to match %q: %v", globPattern, err)
	}

	for _, match := range matches {
		if err := os.RemoveAll(match); err != nil {
			log.Printf("Failed to remove %q: %v", match, err)
		}
	}
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}