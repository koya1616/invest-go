package main

import (
	"api/db"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func fetchHTML(code string) (string, string, error) {
	resp, err := http.Get(fmt.Sprintf("https://finance.yahoo.co.jp/quote/%s.T", code))
	if err != nil {
		return "", "", fmt.Errorf("リクエスト実行エラー: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("レスポンス読み取りエラー: %w", err)
	}
	bodyStr := string(body)

	priceBoard, err := searchHtmlData(bodyStr, "\"mainStocksPriceBoard\":", "otherExchanges")
	if err != nil {
		return "", "", fmt.Errorf("HTML内のデータ読み取りエラー: %w", err)
	}

	price, err := searchHtmlData(priceBoard, "\"price\":\"", "\",")
	if err != nil {
		return "", "", fmt.Errorf("priceデータの読み取りエラー: %w", err)
	}

	dateTime, err := searchHtmlData(priceBoard, "\"priceDateTime\":\"", "\",")
	if err != nil {
		return "", "", fmt.Errorf("priceDateTimeデータの読み取りエラー: %w", err)
	}

	return strings.ReplaceAll(price, ",", ""), dateTime, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func searchHtmlData(str string, startMarker string, endMarker string) (string, error) {
	startIndex := strings.Index(str, startMarker)
	if startIndex == -1 {
		return "", fmt.Errorf("開始マーカーが見つかりません")
	}
	startIndex += len(startMarker)

	endIndex := strings.Index(str[startIndex:], endMarker)
	if endIndex == -1 {
		return "", fmt.Errorf("終了マーカーが見つかりません")
	}

	return str[startIndex : startIndex+endIndex], nil
}

func isWithinTimeRange(now time.Time, start int, end int) bool {
	jstTime := now.In(time.FixedZone("Asia/Tokyo", 9*60*60))
	currentMinutes := jstTime.Hour()*60 + jstTime.Minute()
	return currentMinutes >= start && currentMinutes <= end
}

func main() {
	defer db.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				price, dateTime, err := fetchHTML("7203")
				if err != nil {
					fmt.Printf("HTML取得エラー: %v\n", err)
				}

				now := time.Now()
				if isWithinTimeRange(now, 540, 690) || isWithinTimeRange(now, 750, 930) {
					if strings.Contains(dateTime, "/") {
						if err := db.Instance.InsertTimeSeries("7203", price, fmt.Sprintf("%d/%s 15:30", now.Year(), dateTime)); err != nil {
							fmt.Printf("データ挿入エラー: %v\n", err)
						}
					} else {
						if err := db.Instance.InsertTimeSeries("7203", price, fmt.Sprintf("%s %s", now.Format("2006/01/02"), dateTime)); err != nil {
							fmt.Printf("データ挿入エラー: %v\n", err)
						}
					}
				}
			}
		}
	}()

	http.HandleFunc("/", handler)
	fmt.Println("Server starting on port http://localhost:7778")
	if err := http.ListenAndServe(":7778", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
