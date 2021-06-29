package app

import (
	"fmt"
	"strconv"
	"strings"
)

func AddHandler(chatId int64, words []string) string {
	curChan := make(chan BinanceResponse)
	go GetCourseUsd(words[1], curChan)
	amounts, err := strconv.ParseFloat(words[2], 64)
	if err != nil {
		return "Невозможно определить количество"
	}
	result := <-curChan
	if result.Code != 0 {
		return "Нет такой криптовалюты"
	}

	amount := db.Add(chatId, words[1], amounts)
	return fmt.Sprintf("%f", amount)
}

func SubHandler(chatId int64, words []string) string {
	amounts, err := strconv.ParseFloat(words[2], 64)
	if err != nil {
		return "Невозможно определить количество"
	}
	amount, err := db.Sub(chatId, words[1], amounts)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%f", amount)
}

func DelHandler(chatId int64, currency string) {
	db.Delete(chatId, currency)
}

func ShowHandler(chatId int64, convert string) string {
	currencies := db.GetCurrencies(chatId)
	length := len(currencies)

	var sum float64
	var msg string

	curseChan := make(chan BinanceResponse, length)
	currChan := make(chan BinanceResponse)

	for _, curr := range currencies {
		go GetCourseUsd(curr, curseChan)
	}

	go GetUsdConvert(convert, currChan)
	rub := <-currChan

	fmt.Println(rub)

	for i := 0; i < length; i++ {
		cr := <-curseChan
		curse := strings.ReplaceAll(cr.Symbol, "USDT", "")

		value := db.Ubalances[chatId][curse]
		sum += value * cr.Price * rub.Price

		msg += fmt.Sprintf("%s: %f [%.2f %s]\n", curse, value, cr.Price*rub.Price*value, convert)
	}
	msg += fmt.Sprintf("Total: %.2f %s\n", sum, convert)

	return msg
}
