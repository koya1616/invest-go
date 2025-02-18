package main

import (
	"api/db"
	"api/utils"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if err := db.Instance.DeleteOldTimeSeries("five_minutes_timeseries"); err != nil {
		fmt.Printf("5分足テーブルの今日以前のrecord削除エラー: %v\n", err)
	}
	if err := db.Instance.DeleteOldTimeSeries("one_minute_timeseries"); err != nil {
		fmt.Printf("1分足テーブルの今日以前のrecord削除エラー: %v\n", err)
	}
	if err := db.Instance.DeleteOldTimeSeries("timeseries"); err != nil {
		fmt.Printf("timeseriesテーブルの今日以前のrecord削除エラー: %v\n", err)
	}
	fmt.Fprint(w, "OK")
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
				if now.Second() >= 3 && now.Second() <= 58 {
					continue
				}
				if !utils.IsWithinTimeRange(now, 540, 1020) {
					continue
				}

				code := "7203"
				price, dateTime, err := utils.FetchHTML(code)
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

				if now.Second() != 2 {
					continue
				}

				if err := db.Instance.InsertOneMinuteTimeSeries(); err != nil {
					fmt.Printf("1分足の集計エラー: %v\n", err)
				}

				if err := db.Instance.DeleteDuplicatedTimeSeries("one_minute_timeseries"); err != nil {
					fmt.Printf("1分足の重複削除エラー: %v\n", err)
				}

				if err := db.Instance.InsertFiveMinutesTimeSeries(); err != nil {
					fmt.Printf("5分足の集計エラー: %v\n", err)
				}

				if err := db.Instance.DeleteDuplicatedTimeSeries("five_minutes_timeseries"); err != nil {
					fmt.Printf("5分足の重複削除エラー: %v\n", err)
				}
			}
		}
	}()

	http.HandleFunc("/delete", handleDelete)
	if err := http.ListenAndServe(":7778", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
