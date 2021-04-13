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

	csvPath := getInput("Input CSV Path (absolute path)")
	outputCSVPath := getInput("Directory where Output CSVs will be stored")
	mappingCSVPath := getInput("Mappings CSV Path (absolute path)")

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
		fmt.Printf("error %v", err)
		log.Fatalf("An error occurred while reading input. Please try again.")
		return ""
	}
	input = strings.TrimSuffix(input, "\n")
	if input == "" {
		log.Fatalf("Not a valid input. Please try again\"")
	}
	return input
}

