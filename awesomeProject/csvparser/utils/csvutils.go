package utils

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)


func flushData(writer *csv.Writer) error {
	writer.Flush()
	flushErr := writer.Error() // Checks if any error occurred while writing
	if flushErr != nil {
		fmt.Println("error while writing to the file", flushErr)
		return flushErr
	}
	return nil
}

func closeCSV(csvFile *os.File) error {
	csvCloseError := csvFile.Close()
	if csvCloseError != nil {
		fmt.Println("error while closing the file", csvCloseError)
		return csvCloseError
	}
	return nil
}

func CreateCSV(file string) (bool, error) {
	_, createFileError := os.Create(file)
	if createFileError != nil {
		log.Fatalln("failed to open file", createFileError)
		return false, createFileError
	}
	return true, nil
}

func getCSVForWrite(file string) (*os.File, error) {
	var csvFile *os.File
	path, _ := filepath.Abs(file)
	csvFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if errors.Is(err, os.ErrNotExist) {
		created, error := CreateCSV(file)
		if !created {
			return nil, error
		}
	}
	csvFile, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	return csvFile, nil
}

func WriteDataInCSVSync(file string, data []string) error {
	csvFile, csvOpenError := getCSVForWrite(file)
	if csvOpenError != nil {
		log.Fatal("could not open file to write the data", csvOpenError)
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {}
	}(csvFile)

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	if err := writer.Write(data); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	if flushErr := flushData(writer); flushErr != nil {
		return flushErr
	}

	if closeErr := closeCSV(csvFile); closeErr != nil {
		return closeErr
	}
	return nil
}


func WriteDataInCSV(file string, dataChannel chan []string, processingCompleteChannel chan bool) error {
	csvFile, csvOpenError := getCSVForWrite(file)
	if csvOpenError != nil {
		log.Fatal("could not open file to write the data", csvOpenError)
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {}
	}(csvFile)

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	for data := range dataChannel {
		if err := writer.Write(data); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		if flushErr := flushData(writer); flushErr != nil {
			log.Fatalln("error flushing record to file", flushErr)
			return flushErr
		}
	}

	if closeErr := closeCSV(csvFile); closeErr != nil {
		return closeErr
	}
	processingCompleteChannel <- true
	return nil
}

func ReadDataFromCSV(csvPathString string) [][]string{
	csvPath, _ := filepath.Abs(csvPathString)
	csvFile, openError := os.Open(csvPath)
	if openError != nil {
		log.Fatal("could not open file, error happened ", openError)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	line, _ := reader.ReadAll()

	return line
}