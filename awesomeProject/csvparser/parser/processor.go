package parser

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"csvpaser/utils"
)

func ProcessCSV(csvPathString string, outputCSVPathString string, mappingCSVPathString string) {
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

	firstNameColumn, e1 := FetchHeaderIndex(headers, "First Name")

	lastNameColumn, e1 := FetchHeaderIndex(headers, "Last Name")

	nameColumn, e1 := FetchHeaderIndex(headers, "Name")

	if firstNameColumn == -1 && lastNameColumn == -1 && nameColumn == -1 {
		log.Fatal("name column does not exists")
	}

	emailColumn, e1 := FetchHeaderIndex(headers, "Email")

	if e1 != nil {
		log.Fatal(e1)
	}

	wageColumn, e1 := FetchHeaderIndex(headers, "Wage")
	if e1 != nil {
		log.Fatal(e1)
	}
	employeeNumberColumn, e1 := FetchHeaderIndex(headers, "EmployeeNumber")

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
			utils.WriteDataInCSV(outputCSVPathString + inCorrectOutputCSVName, line)
		} else {
			utils.WriteDataInCSV(outputCSVPathString + correctOutputCSVName, line)
		}

	}
	peopleJson, _ := json.Marshal(people)
	fmt.Println(string(peopleJson))

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
		return strings.Split(record[nameColumnIndex], " ")[1]
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








