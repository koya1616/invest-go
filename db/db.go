package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var Instance *DB

type DB struct {
	*sql.DB
}

func init() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	Instance = &DB{db}
}

func Close() error {
	if Instance != nil {
		return Instance.DB.Close()
	}
	return nil
}

type TimeSeries struct {
	ID       int
	Code     string
	Value    float64
	Datetime time.Time
}

func (db *DB) InsertTimeSeries(code string, value string, date string) error {
	var ts TimeSeries
	err := db.QueryRow(
		"INSERT INTO timeseries (code, value, datetime) VALUES ($1, $2, $3) RETURNING id, code",
		code, value, date,
	).Scan(&ts.ID, &ts.Code)

	if err != nil {
		return fmt.Errorf("error inserting time series: %v", err)
	}

	return nil
}
