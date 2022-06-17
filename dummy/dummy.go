package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var records = readCsvFile("./full_data.csv")

func main() {
	r := gin.Default()
	setupRoutes(r)
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// fmt.Println(records)

}

func setupRoutes(r *gin.Engine) {
	r.GET("/cases/new/country/:country", route1)
	r.GET("/cases/total/country/:from_date", route2)

}

//Dummy function
func route1(c *gin.Context) {
	country, ok := c.Params.Get("country")
	date, date_err := c.GetQuery("date")

	if date_err == false {
		date = "all"
	}
	cases := getNewCases(records, country, date)
	if ok == false {
		res := gin.H{
			"error": "country is missing",
		}
		c.JSON(http.StatusOK, res)
		return
	}

	/*
		city := ""
	*/
	res := gin.H{
		"new_cases": cases,
		"country":   country,
		"date":      date,
	}
	c.JSON(http.StatusOK, res)
}

//Dummy function
func route2(c *gin.Context) {

	date_str, ok := c.Params.Get("from_date")
	date := convertToTimeFormat(date_str)

	total_cases := getTotalCases(records, date)
	if ok == false {
		res := gin.H{
			"error": "date is missing",
			"date":  date,
		}
		c.JSON(http.StatusOK, res)
		return
	}
	/*
		city := ""
	*/
	res := gin.H{
		"total_cases": total_cases,
		"date":        date,
	}
	c.JSON(http.StatusOK, res)
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
func convertToTimeFormat(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		fmt.Println("error: ", err)
	}
	return t
}

func getNewCases(records [][]string, country string, date string) int64 {

	var newCases int64
	var sum int64
	for i := 1; i < len(records); i++ {

		//fmt.Println(records[0][0], i)
		if date != "all" {
			if records[i][1] == country {
				if records[i][0] == date {
					newCases, _ = strconv.ParseInt(records[i][2], 0, 8)
				}

			}
		} else {
			if records[i][1] == country {
				newCases, _ = strconv.ParseInt(records[i][2], 0, 8)
				sum = newCases + sum
			}
		}

	}
	if date != "all" {
		return newCases
	} else {
		return sum
	}

}

func getTotalCases(records [][]string, from_date time.Time) float64 {

	//var total_cases = []int64{}

	var sum = float64(0.0)
	for i := 1; i < len(records); i++ {
		date := convertToTimeFormat(records[i][0])
		//fmt.Println(records[0][0], i)
		if date.After(from_date) {

			//total_cases = append(total_cases, records[i][4])

			temp, err := strconv.ParseFloat(records[i][4], 64)
			if err == nil {
				sum = temp + sum
			}

		}

	}

	return sum
}
