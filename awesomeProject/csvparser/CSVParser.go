package main

import (
	"bufio"
	"csvpaser/parser"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	csvPath := getInput("csvPath")
	outputCSVPath := getInput("outputCSVPath")
	mappingCSVPath := getInput("mappingCSVPath")

	start := time.Now()
	parser.ProcessCSV(csvPath, outputCSVPath, mappingCSVPath)
	elapsed := time.Since(start)
	log.Printf("\nCSVParser took %s", elapsed)
}

func getInput(field string) string {
	fmt.Printf("Enter %s: ", field)
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return ""
	}
	input = strings.TrimSuffix(input, "\n")
	return input
}

