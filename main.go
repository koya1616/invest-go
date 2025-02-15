package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var id int
	var code string
	err = db.QueryRow("SELECT id, code FROM timeseries WHERE id = $1", 1).Scan(&id, &code)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("レコードが見つかりません")
			return
		}
		panic(err)
	}

	fmt.Printf("ID: %d, Code: %s\n", id, code)

	fmt.Println("Server starting on port http://localhost:7778")
	if err := http.ListenAndServe(":7778", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
