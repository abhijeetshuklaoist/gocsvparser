package parser

import (
	"bufio"
	"csvpaser/utils"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func ProcessCSV(csvPathString string, outputCSVPathString string, mappingCSVPathString string, extractData bool) {
	headers := fetchInterestingHeaders(csvPathString, mappingCSVPathString)

	firstNameColumnIndex, _ := FetchHeaderIndex(headers, "First Name")
	lastNameColumnIndex, _ := FetchHeaderIndex(headers, "Last Name")
	nameColumnIndex, _ := FetchHeaderIndex(headers, "Name")

	if firstNameColumnIndex == -1 && lastNameColumnIndex == -1 && nameColumnIndex == -1 {
		log.Fatal("name column does not exists")
	}

	emailColumnIndex, e1 := FetchHeaderIndex(headers, "Email")
	if e1 != nil {
		log.Fatal(e1)
	}

	wageColumnIndex, e1 := FetchHeaderIndex(headers, "Wage")
	if e1 != nil {
		log.Fatal(e1)
	}

	employeeNumberColumnIndex, e1 := FetchHeaderIndex(headers, "EmployeeNumber")
	if e1 != nil {
		log.Fatal(e1)
	}

	timeStamp := time.Now().Format("20060102150405")
	var correctDataChannel = make(chan []string, 500)
	var inCorrectDataChannel = make(chan []string, 500)
	var inCorrectDataProcessingCompleteChannel = make(chan bool)
	var correctDataProcessingCompleteChannel = make(chan bool)

	correctOutputCSVName := path.Base(csvPathString) + "_correct_" + timeStamp
	inCorrectOutputCSVName := path.Base(csvPathString) + "_wrong_" + timeStamp

	go processLine(csvPathString, firstNameColumnIndex, lastNameColumnIndex, nameColumnIndex, emailColumnIndex, wageColumnIndex,
		employeeNumberColumnIndex, correctDataChannel, inCorrectDataChannel, extractData)
	go utils.WriteDataInCSV(outputCSVPathString + inCorrectOutputCSVName, inCorrectDataChannel, inCorrectDataProcessingCompleteChannel)
	go utils.WriteDataInCSV(outputCSVPathString + correctOutputCSVName, correctDataChannel, correctDataProcessingCompleteChannel)

	// wait till correct and incorrect data in written in files
	<-inCorrectDataProcessingCompleteChannel
	<-correctDataProcessingCompleteChannel

	fmt.Printf("\nProcessing completed")
	fmt.Printf("\nInCorrect data file is: %s \n", outputCSVPathString + inCorrectOutputCSVName)
	fmt.Printf("\nCorrect data file is: %s \n", outputCSVPathString + correctOutputCSVName)
}

func fetchInterestingHeaders(csvPathString string, mappingCSVPathString string) map[string]int {
	csvPath, _ := filepath.Abs(csvPathString)
	csvFile, openError := os.Open(csvPath)
	if openError != nil {
		log.Fatal("could not open file, error happened ", openError)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	line, readError := reader.Read()
	headers, err := fetchHeaders(line, fetchInterestingHeaderMappings(mappingCSVPathString))
	if err != nil {
		log.Fatal(readError)
	} else {
		fmt.Println(headers)
	}
	return headers
}

func processLine(csvPathString string, firstNameColumn int, lastNameColumn int, nameColumn int, emailColumn int,
	wageColumn int, employeeNumberColumn int, correctDataChannel chan []string, inCorrectDataChannel chan []string,
	extractData bool) {
	csvPath, _ := filepath.Abs(csvPathString)
	csvFile, _ := os.Open(csvPath)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	headers, _ := reader.Read() // Skip headers
	inCorrectDataChannel <- headers
	if extractData {
		correctDataChannel <- Person{}.fields()
	} else {
		correctDataChannel <- headers
	}

	totalRecordCount := 1
	correctRecordCount := 0
	inCorrectRecordCount := 0
	for {
		line, readError := reader.Read()
		if readError == io.EOF {
			break
		} else if readError != nil {
			log.Fatal(readError)
		}

		person := Person{
			FirstName:      fetchFirstName(line, firstNameColumn, nameColumn),
			LastName:       fetchLastName(line, lastNameColumn, nameColumn),
			Email:          line[emailColumn],
			Wage:           line[wageColumn],
			EmployeeNumber: line[employeeNumberColumn],
		}
		if (person.FirstName == "" && person.LastName == "" )  ||
			person.Email == "" || person.Wage == ""|| person.EmployeeNumber == "" {
			inCorrectDataChannel <- line
			inCorrectRecordCount++
		} else {
			if extractData {
				correctDataChannel <- person.value()
			} else {
				correctDataChannel <- line
			}
			correctRecordCount++
		}
		totalRecordCount++
		if totalRecordCount% 500 == 0 {
			fmt.Printf("\nProcessing.... Processed %d records", totalRecordCount)
		}
	}

	fmt.Printf("\nProcessing completed. Processed %d records", totalRecordCount)
	fmt.Printf("\nProcessing completed. Correct records count is %d", correctRecordCount)
	fmt.Printf("\nProcessing completed. InCorrect records count is %d", inCorrectRecordCount)

	// close channels
	close(inCorrectDataChannel)
	close(correctDataChannel)
}

func fetchHeaders(record []string, headerMappings map[string]string) (map[string]int, error) {
	headers := make(map[string]int)
	for i := 0; i < len(record); i++ {
		expectedHeader, exists := headerMappings[utils.ConvertToLowerCaseString(record[i])]
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
	mappingsCSVPath, _ := filepath.Abs(mappingCSVPath)
	csvFile, err := os.Open(mappingsCSVPath)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, readError := reader.Read()
		if readError == io.EOF {
			break
		} else if readError != nil {
			log.Fatal(readError)
		}
		headers[utils.ConvertToLowerCaseString(line[0])] = utils.ConvertToString(line[1])
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
		nameParts := strings.Split(record[nameColumnIndex], " ")
		if  len(nameParts) == 1{
			return ""
		} else {
			return strings.Split(record[nameColumnIndex], " ")[1]
		}
	}
	log.Fatal("name column not found")
	return ""
}

func FetchHeaderIndex(headers map[string]int, header string) (int, error) {
	headerIndex, exists := headers[header]
	if !exists {
		return -1, fmt.Errorf("column %s does not exists", header)
	} else {
		return headerIndex, nil
	}
}
