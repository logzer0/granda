package main

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func StoreFunctionToDB(functionName, imageName, urlPath string, timeout int) {
	stmt, err := db.Prepare(`INSERT INTO imageFunctions( functionName, imageName, urlPath, timeout  ) VALUES(?,?,?,?)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(functionName, imageName, urlPath, timeout)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}

func StoreRuntimeToDB(functionName, urlPath, invocationTime string, runTimeSecs float64, status string) {
	stmt, err := db.Prepare(`INSERT INTO functionInvocation(functionName, urlPath, invocationTime,runTimeSecs,Status)VALUES(?,?, datetime('now')  ,?,?)`)
	if err != nil {
		panic(err)
	}
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(functionName, urlPath, runTimeSecs, status)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}

func getFunctionsFromDB() {
	var functionName, imageName, urlPath string
	var timeout int
	rows, err := db.Query(`SELECT functionName, imageName, urlPath, timeout FROM imageFunctions`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&functionName, &imageName, &urlPath, &timeout)
		if err != nil {
			panic(err)
		}
		pathToImage[urlPath] = Image{funcName: functionName, imageName: imageName, Timeout: timeout}
	}
}

func getInvocationCount(duration int) (count int) {
	var functionName string
	var total, year, month, date, hours, minutes, seconds, diff int
	rows, err := db.Query(`select functionName, count(runTimeSecs) as total ,  strftime('%Y', invocationTime) as year, strftime('%m', invocationTime) as month, strftime('%d', invocationTime) as date,
	strftime('%H', invocationTime) as hours, strftime('%M', invocationTime) as minutes, strftime('%S',invocationTime) as seconds, 
	strftime('%s','now') - strftime('%s',invocationTime) as diff
	from functionInvocation where diff < ?
	group by year, month,date,hours,minutes, seconds/?`, duration, duration)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&functionName, &total, &year, &month, &date, &hours, &minutes, &seconds, &diff)
		if err != nil {
			panic(err)
		}
		count = total
	}
	return
}

func getInvocationDetailsFromDB(functionName string) (details [][]string) {
	var invocationTime, status string
	var runTimeSecs float64
	rows, err := db.Query(`SELECT invocationTime, runTimeSecs, Status FROM functionInvocation  order by invocationTime desc`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&invocationTime, &runTimeSecs, &status)
		if err != nil {
			panic(err)
		}
		details = append(details, []string{functionName, invocationTime, strconv.FormatFloat(runTimeSecs, 'f', -1, 64), status})
	}
	return
}
