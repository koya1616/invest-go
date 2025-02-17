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
	currentMinutes := now.Hour()*60 + now.Minute()
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
				now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
				if now.Second() != 59 && now.Second() != 58 && now.Second() != 00 && now.Second() != 1 {
					continue
				}
				if !isWithinTimeRange(now, 540, 690) && !isWithinTimeRange(now, 750, 930) {
					continue
				}

				code := "7203"
				price, dateTime, err := fetchHTML(code)
				if err != nil {
					fmt.Printf("HTML取得エラー: %v\n", err)
					continue
				}

				fmt.Println(now.Format("2006-01-02 15:04:05"), price, dateTime)

				var formattedDateTime string
				if strings.Contains(dateTime, "/") {
					formattedDateTime = fmt.Sprintf("%d/%s 15:30", now.Year(), dateTime)
				} else {
					formattedDateTime = fmt.Sprintf("%s %s", now.Format("2006/01/02"), dateTime)
				}

				if err := db.Instance.InsertTimeSeries(code, price, formattedDateTime); err != nil {
					fmt.Printf("データ挿入エラー: %v\n", err)
					continue
				}
			}
		}
	}()

	if err := http.ListenAndServe(":7778", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
