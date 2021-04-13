package main

import (
	"csvpaser/utils"
	"log"
	"strconv"
	"time"
)

func createData()  {
	start := time.Now()
	csvPathString := "csvs/data/rostertest.csv"
	var data = []string {"Name" , "Email", "ExtraColumn1", "ExtraColumn2", "ExtraColumn3", "Salary", "ID", "ExtraColumn4",
		"ExtraColumn5", "ExtraColumn6", "ExtraColumn7", "ExtraColumn8", "ExtraColumn9"}
	utils.WriteDataInCSVSync(csvPathString, data)
	for i:= 1; i < 1000000; i ++ {
		var data = []string {"Name" + strconv.Itoa(i), "EmailSomething@gmai.com"+ strconv.Itoa(i),
			"ExtraColumn1" +  strconv.Itoa(i), "ExtraColumn2" +  strconv.Itoa(i), "ExtraColumn3" +  strconv.Itoa(i),
			"Salary" + strconv.Itoa(i), "ID" + strconv.Itoa(i),
			"ExtraColumn4" +  strconv.Itoa(i), "ExtraColumn5" +  strconv.Itoa(i), "ExtraColumn6" +  strconv.Itoa(i),
			"ExtraColumn7" +  strconv.Itoa(i), "ExtraColumn8" +  strconv.Itoa(i), "ExtraColumn9" +  strconv.Itoa(i)}
		utils.WriteDataInCSVSync(csvPathString, data)

		var wrongdata = []string {"Name" + strconv.Itoa(i), "",
			"ExtraColumn1" +  strconv.Itoa(i), "ExtraColumn2" +  strconv.Itoa(i), "ExtraColumn3" +  strconv.Itoa(i),
			"Salary" + strconv.Itoa(i), "ID" + strconv.Itoa(i),
			"ExtraColumn4" +  strconv.Itoa(i), "ExtraColumn5" +  strconv.Itoa(i), "ExtraColumn6" +  strconv.Itoa(i),
			"ExtraColumn7" +  strconv.Itoa(i), "ExtraColumn8" +  strconv.Itoa(i), "ExtraColumn9" +  strconv.Itoa(i)}
		utils.WriteDataInCSVSync(csvPathString, wrongdata)
	}
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

func main1() {
	createData()
}
