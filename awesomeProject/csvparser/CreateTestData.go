package main

import (
	"csvpaser/utils"
	"strconv"
)

func createData()  {
	csvPathString := "csvs/data/rostertest.csv"
	var data = []string {"Name" , "Email", "Salary", "ID"}
	utils.WriteDataInCSVSync(csvPathString, data)
	for i:= 1; i < 1000000; i++ {
		var data = []string {"Name" + strconv.Itoa(i), "EmailSomething@gmai.com"+ strconv.Itoa(i),
			"Salary" + strconv.Itoa(i), "ID" + strconv.Itoa(i)}
		utils.WriteDataInCSVSync(csvPathString, data)
	}

}

func main1() {
	createData()
}
