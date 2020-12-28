package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
	"strings"
)

type Event struct {
	Code    string
	OldDate string
	Venue 	string
	NewDate string
	NewCode string
}

type CurrentEvent struct {
	Code      string
	StartDate string
	EndDate	  string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "test"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	// Open the file
	csvFile, _ := os.Open("/Users/leoclarke/Desktop/Sheet1.csv")

	reader := csv.NewReader(bufio.NewReader(csvFile))
	var events []Event
	for {
		col, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		events = append(events, Event{
			Code: col[0],
			OldDate:  col[1],
			Venue:  col[2],
			NewDate:  col[3],
			NewCode: generateNewCode(col[0],col[3]),
		})
	}

	for _, x := range events[1:]{
		fmt.Printf("%+v\n", x)
	}
	generateQuery(events[2].Code, events[2].NewCode, events[2].NewDate)
}


func generateNewCode(oldCode, newDate string) string {
	date := ""
	day := getDay(newDate)
	code := trimDate(oldCode, 6)
	if strings.Contains(newDate, "March"){
		date = "2103"
	}else if strings.Contains(newDate, "April") {
		date = "2104"
	}else{
		date = ""
	}
	var newString = date + day + code
	return newString
}

func generateQuery(oldCode, newCode, newDate string) string {
	var ce CurrentEvent

	db := dbConn()

	err := db.QueryRow("SELECT start_date, end_date FROM events WHERE code = ? ", oldCode).Scan(&ce.StartDate, &ce.EndDate)
	if err != nil {
		log.Println("Event not found")
	} else {
		newStartDate := generateDateTime(ce.StartDate, newDate)
		newEndDate := generateDateTime(ce.EndDate, newDate)
		insForm, err := db.Prepare("UPDATE events SET code = ?, start_date = ?, end_date = ? WHERE code = ?")
		if err != nil {
			panic(err.Error())
		}
		_, _ = insForm.Exec(newCode, newStartDate, newEndDate, oldCode)
		fmt.Print("UPDATE events SET code =", "'" ,newCode, "'" ,", start_date =","'" ,newStartDate,"'" , ", end_date = ","'" , newEndDate,"'" , ",WHERE code =","'" ,oldCode,"'" )
	}

	query := "update "
	return query
}

func generateDateTime(dateTime, newDate string) string {
	time := trimDate(dateTime, 11)
	year := "2021"
	month := getMonth(newDate)
	day := getDay(newDate)
	newDateTime := year + "-" + month + "-" + day + " " + time

	return newDateTime
}

func getMonth(newDate string) string  {
	month := ""
	if strings.Contains(newDate, "March"){
		month = "03"
	}else if strings.Contains(newDate, "April") {
		month = "04"
	}else if strings.Contains(newDate, "May") {
		month = "05"
	}

	return month
}

func trimDate(s string, l int) string {
	m := 0
	for i := range s {
		if m >= l {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

func getDay(date string) string  {
	day := ""
	if strings.Contains(date, "1"){
		day = "01"
	}else if strings.Contains(date, "2") {
		day = "02"
	}else if strings.Contains(date, "3") {
		day = "03"
	}
	return day
}