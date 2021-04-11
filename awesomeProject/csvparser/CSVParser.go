package main

import (
	"csvpaser/parser"
)

func main() {
	csvPath := "csvs/data/roster1.csv"
	outputCSVPath := "csvs/output/"
	mappingCSVPath := "utils/mappings.csv"
	parser.ProcessCSV(csvPath, outputCSVPath, mappingCSVPath)
}

