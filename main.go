package main

import (
	"api/db"
	"fmt"
	"io"
	"net/http"
)

func fetchHTML() (string, error) {
	resp, err := http.Get("https://finance.yahoo.co.jp/quote/7203.T")
	if err != nil {
		return "", fmt.Errorf("リクエスト実行エラー: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("レスポンス読み取りエラー: %w", err)
	}

	return string(body), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	html, err := fetchHTML()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		return
	}
	fmt.Println(html)

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
