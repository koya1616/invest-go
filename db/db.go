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

func (db *DB) InsertOneMinuteTimeSeries() error {
	query := `
		INSERT INTO one_minute_timeseries (code, value, datetime)
		WITH latest_records AS (
			SELECT code, datetime, MAX(id) as latest_id
			FROM timeseries
			GROUP BY code, datetime
		)
		SELECT t.code, t.value, t.datetime
		FROM timeseries t
		INNER JOIN latest_records l ON t.id = l.latest_id;
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error executing query to insert into one_minute_timeseries: %v", err)
	}

	return nil
}

func (db *DB) InsertFiveMinutesTimeSeries() error {
	query := `
		WITH five_minutes_groups AS (
			SELECT
				code,
				date_trunc('hour', datetime) + INTERVAL '5 minutes' * (date_part('minute', datetime)::integer / 5) AS five_min_datetime,
				MAX(id) as max_id
			FROM one_minute_timeseries
			GROUP BY
				code,
				date_trunc('hour', datetime) + INTERVAL '5 minutes' * (date_part('minute', datetime)::integer / 5)
		)
		INSERT INTO five_minutes_timeseries (code, value, datetime)
		SELECT g.code, t.value, g.five_min_datetime
		FROM five_minutes_groups g
		JOIN one_minute_timeseries t ON t.id = g.max_id
		ORDER BY g.five_min_datetime;
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error executing query to insert into five_minutes_timeseries: %v", err)
	}

	return nil
}

func (db *DB) DeleteDuplicatedTimeSeries(table string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id NOT IN (
			SELECT MAX(id)
			FROM %s
			GROUP BY code, datetime
		);
	`, table, table)

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error executing query to delete from %s: %v", table, err)
	}

	return nil
}

func (db *DB) DeleteOldTimeSeries(table string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE datetime < CURRENT_DATE;
	`, table)

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error executing query to delete old %s records: %v", table, err)
	}

	return nil
}
