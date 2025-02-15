package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
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
	ID   int
	Code string
}

func (db *DB) GetTimeSeriesById(id int) (*TimeSeries, error) {
	var ts TimeSeries
	err := db.QueryRow("SELECT id, code FROM timeseries WHERE id = $1", id).Scan(&ts.ID, &ts.Code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying time series: %v", err)
	}
	return &ts, nil
}
