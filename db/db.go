package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
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
