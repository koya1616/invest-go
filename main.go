package main

import (
	"api/db"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
	defer db.Close()

	ts, err := db.Instance.GetTimeSeriesById(1)
	if err != nil {
		panic(err)
	}
	if ts == nil {
		fmt.Println("レコードが見つかりません")
		return
	}

	fmt.Printf("ID: %d, Code: %s\n", ts.ID, ts.Code)

	http.HandleFunc("/", handler)
	fmt.Println("Server starting on port http://localhost:7778")
	if err := http.ListenAndServe(":7778", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
