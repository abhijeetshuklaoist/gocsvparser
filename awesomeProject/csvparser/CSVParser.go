package main

import (
	"csvpaser/parser"
	"log"
	"time"
)

func main() {
	start := time.Now()
	csvPath := "csvs/data/rostertest.csv"
	outputCSVPath := "csvs/output/"
	mappingCSVPath := "utils/mappings.csv"
	parser.ProcessCSV(csvPath, outputCSVPath, mappingCSVPath)
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

