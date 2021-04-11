package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	csvPath := "csvs/roster1.csv"
	outputCSVPath := "csvs/"
	mappingCSVPath := "utils/mappings.csv"
	processCSV(csvPath, outputCSVPath, mappingCSVPath)
}

func processCSV(csvPathString string, outputCSVPathString string, mappingCSVPathString string) {
	csvPath, _ := filepath.Abs(csvPathString)
	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Printf("error happend %s", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var people []Person
	line, error := reader.Read()
	headers, err := fetchHeaders(line, fetchInterestingHeaderMappings(mappingCSVPathString))
	if err != nil {
		log.Fatal(error)
	} else {
		fmt.Println(headers)
	}

	firstNameColumn, e1 := fetchHeaderIndex(headers, "First Name")

	lastNameColumn, e1 := fetchHeaderIndex(headers, "Last Name")

	nameColumn, e1 := fetchHeaderIndex(headers, "Name")

	if firstNameColumn == -1 && lastNameColumn == -1 && nameColumn == -1 {
		log.Fatal("name column does not exists")
	}

	emailColumn, e1 := fetchHeaderIndex(headers, "Email")

	if e1 != nil {
		log.Fatal(e1)
	}

	wageColumn, e1 := fetchHeaderIndex(headers, "Wage")
	if e1 != nil {
		log.Fatal(e1)
	}
	employeeNumberColumn, e1 := fetchHeaderIndex(headers, "EmployeeNumber")

	if e1 != nil {
		log.Fatal(e1)
	}


	timeStamp := time.Now().Format("20060102150405")
	correctOutputCSVName := path.Base(csvPathString) + "_correct_" + timeStamp
	inCorrectOutputCSVName := path.Base(csvPathString) + "_wrong_" + timeStamp

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		person := Person{
			FirstName:      fetchFirstName(line, firstNameColumn, nameColumn),
			LastName:       fetchLastName(line, lastNameColumn, nameColumn),
			Email:          line[emailColumn],
			Wage:           line[wageColumn],
			EmployeeNumber: line[employeeNumberColumn],
		}
		people = append(people, person)
		if (person.FirstName == "" && person.LastName == "" )  ||
					person.Email == "" || person.Wage == ""|| person.EmployeeNumber == "" {
			writeDataInCSV(outputCSVPathString + inCorrectOutputCSVName, line)
		} else {
			writeDataInCSV(outputCSVPathString + correctOutputCSVName, line)
		}

	}
	peopleJson, _ := json.Marshal(people)
	fmt.Println(string(peopleJson))

}

func fetchHeaders(record []string, headerMappings map[string]string) (map[string]int, error) {
	headers := make(map[string]int)
	for i := 0; i < len(record); i++ {
		expectedHeader, exists := headerMappings[convertToLowerCaseString(record[i])]
		if !exists {
			continue // Current header is not the interesting to us
		}
		header, exists := headers[expectedHeader]
		if exists {
			return nil, fmt.Errorf("duplicate column %d exists", header)
		} else {
			ByteOrderMarkAsString := string('\uFEFF')
			headerString := strings.TrimPrefix(expectedHeader, ByteOrderMarkAsString)
			headers[headerString] = i
		}
	}
	return headers, nil
}

func fetchInterestingHeaderMappings(mappingCSVPath string) map[string]string {
	headers := make(map[string]string)
	path, _ := filepath.Abs(mappingCSVPath)
	csvFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		headers[convertToLowerCaseString(line[0])] = convertToString(line[1])
	}

	return headers
}

func fetchFirstName(record []string, firstNameColumnIndex int, nameColumnIndex int) string {
	if firstNameColumnIndex != -1 {
		return record[firstNameColumnIndex]
	}
	if nameColumnIndex != -1 {
		return strings.Split(record[nameColumnIndex], " ")[0]
	}
	log.Fatal("name column not found")
	return ""
}

func fetchLastName(record []string, lastNameColumnIndex int, nameColumnIndex int) string {
	if lastNameColumnIndex != -1 {
		return record[lastNameColumnIndex]
	}
	if nameColumnIndex != -1 {
		return strings.Split(record[nameColumnIndex], " ")[1]
	}
	log.Fatal("name column not found")
	return ""
}

func fetchHeaderIndex(headers map[string]int, header string) (int, error) {
	headerIndex, exists := headers[header]
	if !exists {
		return -1, fmt.Errorf("column %s does not exists", header)
	} else {
		return headerIndex, nil
	}
}

func writeDataInCSV(file string, data []string) error {
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

type Person struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Wage           string `json:"wage"`
	EmployeeNumber string `json:"number"`
}

func convertToString(s string) string {
	ByteOrderMarkAsString := string('\uFEFF')
	str := strings.TrimPrefix(s, ByteOrderMarkAsString)
	return str
}

func convertToLowerCaseString(s string) string {
	str := strings.ToLower(convertToString(s))
	return str
}

func createCSV(file string) (bool, error) {
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
		created, error := createCSV(file)
		if !created {
			return nil, error
		}
	}
	csvFile, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	return csvFile, nil
}

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
