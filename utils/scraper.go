package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

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

func FetchHTML(code string) (string, string, error) {
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

func FetchToken(code string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://finance.yahoo.co.jp/quote/%s.T", code))
	if err != nil {
		return "", fmt.Errorf("リクエスト実行エラー: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("レスポンス読み取りエラー: %w", err)
	}
	bodyStr := string(body)

	token, err := searchHtmlData(bodyStr, "\"stocksJwtToken\":\"", "\",")
	if err != nil {
		return "", fmt.Errorf("HTML内のデータ読み取りエラー: %w", err)
	}

	return token, nil
}
