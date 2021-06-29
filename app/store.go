package app

import (
	"errors"
	"sync"
)

var mutex sync.Mutex

type Wallet map[string]float64

type Db struct {
	Ubalances map[int64]Wallet
}

var db Db

func init() {
	db = Db{
		Ubalances: map[int64]Wallet{},
	}
}

func GetDb() Db {
	return db
}

func (db *Db) Add(chatId int64, curr string, value float64) float64 {
	mutex.Lock()
	_, ok := db.Ubalances[chatId]
	if !ok {
		db.Ubalances[chatId] = Wallet{}
	}
	db.Ubalances[chatId][curr] += value
	mutex.Unlock()

	return db.Ubalances[chatId][curr]
}

func (db *Db) Sub(chatId int64, curr string, value float64) (float64, error) {
	mutex.Lock()
	_, ok := db.Ubalances[chatId]
	if !ok {
		mutex.Unlock()
		return 0, errors.New("Нет такой валюты")
	}
	if db.Ubalances[chatId][curr] < value {
		mutex.Unlock()
		return 0, errors.New("Недостаточно номинала на счету")
	}
	db.Ubalances[chatId][curr] -= value
	mutex.Unlock()

	return db.Ubalances[chatId][curr], nil
}

func (db *Db) Delete(chatId int64, curr string) {
	mutex.Lock()
	delete(db.Ubalances[chatId], curr)
	mutex.Unlock()
}

func (db *Db) GetCurrencies(chatId int64) []string {
	currencies := []string{}
	for key := range db.Ubalances[chatId] {
		currencies = append(currencies, key)
	}
	return currencies
}
