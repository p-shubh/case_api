package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

/*1. GET /cases/new/country/:country?date=2020-12-18

3. GET /cases/total/country/:from_date */

// var records = readCsvFile("./full_data.csv")

func main() {
	// records := readCsvFile("./full_data.csv")

	r := gin.Default()
	setupRoutes(r)
	r.Run(":8082") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// fmt.Println(records)
}

func setupRoutes(r *gin.Engine) {
	r.GET("/cases/new/country/:country", setRoute1)
	r.GET("/cases/total/country/:from_date", setRoute2)

}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

//Dummy function
func setRoute1(c *gin.Context) {

	records := readCsvFile("./full_data.csv")
	country, ok := c.Params.Get("country")
	date, ok := c.GetQuery("date")

	cases := getNewCaseStatus(country, date, records)

	if ok == false {
		res := gin.H{
			"error": "file_is_missing",
		}
		c.JSON(http.StatusOK, res)
		return
	}
	/*
		city := ""
	*/
	res := gin.H{ //response
		"new_case": cases,
		"date":     date,
		"country":  country,
		"count":    len(cases),
	}
	c.JSON(http.StatusOK, res)
}

func getNewCaseStatus(country string, date string, records [][]string) []string {

	var new_cases []string
	for i := 0; i < len(records); i++ {
		if records[i][0] == date && records[i][1] == country {
			// new_cases = records[i][2]
			new_cases = append(new_cases, records[i][2])
			break
		}
	}
	return new_cases
}

func setRoute2(c *gin.Context) {
	// country, ok := c.Params.Get("country")
	date, ok := c.GetQuery("date")

	records := readCsvFile("./full_data.csv")
	total_cases := getTotalCasesStatus(date, records)

	if ok == false {
		res := gin.H{
			"error": "file_is_missing",
			"date":  date,
		}
		c.JSON(http.StatusOK, res)
		return
	}
	/*
		all_cases := ""
	*/
	res := gin.H{ //response
		"total_cases": total_cases,
		"date":        date,
		// "country":     country,
		"count": len(total_cases),
	}
	c.JSON(http.StatusOK, res)
}

func getTotalCasesStatus(date string, records [][]string) []string {
	var all_cases []string

	for i := 0; i < len(records); i++ {
		if records[i][0] == date {
			all_cases = append(all_cases, records[i][4])
		}
	}
	return all_cases
}

// var full_data struct {
// 	date            string
// 	location        string
// 	new_cases       float64
// 	new_deaths      float64
// 	total_cases     float64
// 	total_deaths    float64
// 	weekly_cases    float64
// 	weekly_deaths   float64
// 	biweekly_cases  float64
// 	biweekly_deaths float64
// }
