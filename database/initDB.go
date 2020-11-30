package database

import (
	"database/sql"
	"fmt"
	"log"
	"teamAPI/config"
)

var DB *sql.DB

func OpenDB() *sql.DB {
	db, err := sql.Open("mysql", config.DatabaseUsername+":"+config.DatabasePassword+"@tcp("+config.DatabaseHost+")/"+config.DatabaseName)
	if err != nil {
		log.Println(err.Error())
	}

	return db
}

func CreateDatabaseAndTables() {
	config.DB = OpenDB()

	_, err := config.DB.Query(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.DatabaseName) )
	if err != nil {
		log.Fatal("ERROR: Creating database:", err.Error())
		return
	}
	//_, err = db.DB.Exec(fmt.Sprintf("USE %s", config.DatabaseName))
	//_, err = db.DB.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id int NOT NULL AUTO_INCREMENT, symbol VARCHAR(10), price float, PRIMARY KEY (id))", config.PriceDataTable))
	//if err != nil {
	//	log.Fatal("ERROR: Creating "+config.PriceDataTable+" table:", err.Error())
	//	return
	//}
	//candlestickTable, err := db.DB.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id int NOT NULL AUTO_INCREMENT, symbol VARCHAR(10), open float, high float, low float, close float, PRIMARY KEY (id))", config.CandlestickTable))
	//defer candlestickTable.Close()
	//if err != nil {
	//	log.Fatal("ERROR: Creating "+config.CandlestickTable+"table:", err.Error())
	//	return
	//}
	return
}