package app

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
)

type BinanceResponse struct {
    Price  float64 `json:"price,string"`
    Code   int64   `json:"code"`
    Symbol string  `json:"symbol"`
}

func GetCourseUsd(currency string, curseChan chan BinanceResponse) (jsonResponse BinanceResponse, err error) {
    jsonResponse, err = makeRequest(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", currency))

    curseChan <- jsonResponse
    return
}

func GetUsdConvert(currency string, rubChan chan BinanceResponse) (err error) {
    defer close(rubChan)
    jsonResponse, err := makeRequest(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=USDT%s", currency))

    rubChan <- jsonResponse
    return
}

func makeRequest(url string) (jsonResponse BinanceResponse, err error) {
    resp, err := http.Get(url)
    if err != nil {
        return
    }

    defer resp.Body.Close()

    err = json.NewDecoder(resp.Body).Decode(&jsonResponse)

    if err != nil {
        return
    }
    if jsonResponse.Code != 0 {
        err = errors.New("Неверный символ")
    }
    return
}
