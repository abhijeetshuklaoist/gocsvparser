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
)

func main() {

	var x int                            // variable declaration
	var y = 11                           // Inline initialization
	x = 5                                // initialization
	arr := [2][2]int{{1, 2}, {3, x + y}} // Fix sized 2D array declaration and inline initialization
	strings := make([]string, 1)         //Make a slice
	fmt.Printf("\nSize of Slice is %v and capacity of slice is %v", len(strings), cap(strings))
	strings = append(strings, "One", "Two", "Three", "Four") //Append values to a slice
	fmt.Printf("\nHello World. Arrays is %v and Slice is %v", arr, strings)
	fmt.Printf("\nSize of Slice is %v and capacity of slice is %v", len(strings), cap(strings))
	readCSV()
	customCSV()
}

func readCSV() {
	path, _ := filepath.Abs("csvparser/csvs/roster1.csv")
	csvFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error happend", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var people []Person
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		people = append(people, Person{
			Name:   line[0],
			Email:  line[1],
			Wage:   line[2],
			Number: convertToInt(line[3]),
		})
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

func convertToInt(s string) int {
	flt, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return int(flt)
	}
	return 0
}

type Person struct {
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Wage   string `json:"wage"`
	Number int     `json:"number""`
}

func customCSV() {
	path, _ := filepath.Abs("csvparser/csvs/roster1.csv")
	csvFile, err := os.Open(path)
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