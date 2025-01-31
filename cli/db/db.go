package db

import (
	"UptimeKumaProbe/helpers"
	"database/sql"
	"errors"
	"os"
	"time"
	_ "modernc.org/sqlite"
)

//const dbPath = "/opt/kprobe/db/db.sqlite"
const dbPath = "db.sqlite" //TODO temporary, change

var DB *sql.DB

func connectDatabase() {
	if !DatabaseExist() {
		helpers.PrintError(true, "Database does not exist, run <kprobe db init> first")
	}

	var err error

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		helpers.PrintError(true, "Failed to connect to database (" + err.Error() + ")")
	}
}

func InitDatabase() {
	var err error

	if DatabaseExist() {
		err = os.Remove(dbPath)
		if err != nil {
			helpers.PrintError(true, "Failed to delete database (" + err.Error() + ")")
		}
	}

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		helpers.PrintError(true, "Failed to connect to database (" + err.Error() + ")")
	}

	createTableQuery := `
	CREATE TABLE keys (
		name VARCHAR(32) PRIMARY KEY,
		value VARCHAR(4096) NOT NULL
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to create table (" + err.Error() + ")")
	}

	createTableQuery = `
	CREATE TABLE scans (
		name VARCHAR(32) PRIMARY KEY,
		type VARCHAR(4) CHECK(type IN ('ping', 'http')),
		address VARCHAR(256) NOT NULL,
		timeout INTEGER,
		status_code VARCHAR(256),
		keyword TEXT
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to create table (" + err.Error() + ")")
	}

	createTableQuery = `
	CREATE TABLE history (
		scan_name VARCHAR(32) NOT NULL,
		generated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		passed BOOLEAN NOT NULL,
		delete_after DATETIME NOT NULL,
		PRIMARY KEY (scan_name, generated)
	);`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to create table (" + err.Error() + ")")
	}

	InsertValue("probe_name", "New Probe")
	InsertValue("db_version", "v1.0")
	InsertValue("db_init_time", time.Now().String()) 
	
	InsertValue("config_set", "false")
	InsertValue("delete_after", "7")

	InsertValue("api_port", "80")
	InsertValue("editor_endpoint", "true")

	InsertValue("ping_retries", "5")
}

func DatabaseExist() bool {
	if _, err := os.Stat(dbPath); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		helpers.PrintError(true, "Failed to check database existence (" + err.Error() + ")")
	}

	return false
}

func GetValue(key string) string {
	if DB == nil {
		connectDatabase()
	}

	var value string

	query := `
	SELECT value
	FROM keys
	WHERE name = ?;
	`

	err := DB.QueryRow(query, key).Scan(&value)
	if err != nil {
		helpers.PrintError(true, "Failed to get data from database (" + err.Error() + ")")
	}

	return value
}

func InsertValue(key string, value string) {
	if DB == nil {
		connectDatabase()
	}

	insertQuery := `
	INSERT INTO keys (name, value) 
	VALUES (?, ?) 
	ON CONFLICT(name) 
	DO UPDATE SET value = excluded.value;
	`

	_, err := DB.Exec(insertQuery, key, value)
	if err != nil {
		helpers.PrintError(true, "Failed to insert data into database (" + err.Error() + ")")
	}
}

func GetScans() []helpers.Scan {
	if DB == nil {
		connectDatabase()
	}

	var scans []helpers.Scan

	query := `
	SELECT name, type, address, timeout, status_code, keyword
	FROM scans;
	`

	rows, err := DB.Query(query)
	if err != nil {
		helpers.PrintError(true, "Failed to get data from database (" + err.Error() + ")")
	}

	for rows.Next() {
		var scan helpers.Scan

		err = rows.Scan(&scan.Name, &scan.Type, &scan.Address, &scan.Timeout, &scan.StatusCode, &scan.Keyword)
		if err != nil {
			helpers.PrintError(true, "Failed to scan data from database (" + err.Error() + ")")
		}

		scans = append(scans, scan)
	}

	return scans
}

func AddScan(scan helpers.Scan) {
	if DB == nil {
		connectDatabase()
	}

	insertQuery := `
	INSERT INTO scans (name, type, address, timeout, status_code, keyword) 
	VALUES (?, ?, ?, ?, ?, ?);
	`

	_, err := DB.Exec(insertQuery, scan.Name, scan.Type, scan.Address, scan.Timeout, scan.StatusCode, scan.Keyword)
	if err != nil {
		helpers.PrintError(true, "Failed to insert data into database (" + err.Error() + ")")
	}
}


func DeleteScans() {
	if DB == nil {
		connectDatabase()
	}

	deleteQuery := `
	DELETE FROM scans;
	`

	_, err := DB.Exec(deleteQuery)
	if err != nil {
		helpers.PrintError(true, "Failed to delete scans from database (" + err.Error() + ")")
	}
}