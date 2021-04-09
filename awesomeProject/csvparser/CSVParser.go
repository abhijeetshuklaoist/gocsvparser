package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	/*
		var x int                            // variable declaration
		var y = 11                           // Inline initialization
		x = 5                                // initialization
		arr := [2][2]int{{1, 2}, {3, x + y}} // Fix sized 2D array declaration and inline initialization
		strings := make([]string, 1)         //Make a slice
		fmt.Printf("\nSize of Slice is %v and capacity of slice is %v", len(strings), cap(strings))
		strings = append(strings, "One", "Two", "Three", "Four") //Append values to a slice
		fmt.Printf("\nHello World. Arrays is %v and Slice is %v", arr, strings)
		fmt.Println("\nSize of Slice is %v and capacity of slice is %v", len(strings), cap(strings))
	*/

	readCSV()
	//customCSV()
}

func readCSV() {
	path, _ := filepath.Abs("csvparser/csvs/roster1.csv")
	csvFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error happend", err)
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
		people = append(people, Person{
			FirstName:      fetchFirstName(line, firstNameColumn, nameColumn),
			LastName:      fetchLastName(line, lastNameColumn, nameColumn),
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
		expectedHeader, exists := headerMappings[convertToString(record[i])]
		if !exists {
			continue // Current header is not the interesting to us
		}
		header, exists := headers[expectedHeader]
		if exists {
			return nil, fmt.Errorf("duplicate column %s exists", header)
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
	path, _ := filepath.Abs("csvparser/utils/mappings.csv")
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
		headers[convertToString(line[0])] = convertToString(line[1])
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

type Person struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Wage           string `json:"wage"`
	EmployeeNumber string `json:"number"`
}

func customCSV() {
	path, _ := filepath.Abs("csvparser/csvs/roster1.csv")
	csvFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	people := []*Person{}
	if err := gocsv.UnmarshalFile(csvFile, &people); err != nil { // Load clients from file
		panic(err)
	}
	peopleJson, _ := json.Marshal(people)
	fmt.Println(string(peopleJson))
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


func convertToInt(s string) int {
	flt, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return int(flt)
	}
	return 0
}
