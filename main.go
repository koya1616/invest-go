package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)

	connStr := "postgres://postgres:password@db:5432/db?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Server starting on port http://localhost:7778")
	if err := http.ListenAndServe(":7778", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
