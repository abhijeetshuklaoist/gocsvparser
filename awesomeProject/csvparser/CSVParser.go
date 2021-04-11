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
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	readCSV()
}

func readCSV() {
	path, _ := filepath.Abs("csvs/roster1.csv")
	csvFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("error happend %s", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var people []Person
	line, error := reader.Read()
	headers, err := fetchHeaders(line, fetchInterestingHeaderMappings())
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

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		writeCorrectDataInCSV("csvs/output1.csv", line)

		people = append(people, Person{
			FirstName:      fetchFirstName(line, firstNameColumn, nameColumn),
			LastName:       fetchLastName(line, lastNameColumn, nameColumn),
			Email:          line[emailColumn],
			Wage:           line[wageColumn],
			EmployeeNumber: line[employeeNumberColumn],
		})
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

func fetchInterestingHeaderMappings() map[string]string {
	headers := make(map[string]string)
	path, _ := filepath.Abs("utils/mappings.csv")
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

func writeCorrectDataInCSV(file string, data []string) error {
	csvFile, error := getCSVForWrite(file)
	if error != nil {
		log.Fatal("could not open file to write the data", error)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	if err := writer.Write(data); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	writer.Flush()
	err := writer.Error() // Checks if any error occurred while writing
	if err != nil {
		fmt.Println("Error while writing to the file ::", err)
		return err
	}
	csvFile.Close()
	return nil
}

type Person struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Wage           string `json:"wage"`
	EmployeeNumber string `json:"number"`
}

func convertToFloat(s string) float64 {
	flt, err := strconv.ParseFloat(s, 32)
	if err == nil {
		return flt
	}
	return 0
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

func convertToInt(s string) int {
	flt, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return int(flt)
	}
	return 0
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
