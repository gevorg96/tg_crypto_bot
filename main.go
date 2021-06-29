package main

import (
	"log"
	"skillbox/app"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("BOT_APIKEY_HERE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		go handleCommand(bot, update)

	}
}

/*type wallet map[string]float64

var db = map[int64]wallet{}*/

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	commands := strings.Split(update.Message.Text, " ")
	words := toUpperSlice(commands)
	chatId := update.Message.Chat.ID
	var msg string

	switch words[0] {
	case "ADD":
		if validateCommand(words, 3) {
			msg = app.AddHandler(chatId, words)
			break
		}
		msg = "Неверная команда"

	case "SUB":
		if validateCommand(words, 3) {
			msg = app.SubHandler(chatId, words)
			break
		}
		msg = "Неверная команда"

	case "DEL":
		if validateCommand(words, 2) {
			app.DelHandler(chatId, words[1])
			break
		}
		msg = "Неверная команда"

	case "SHOW":
		msg = app.ShowHandler(chatId, "RUB")

	default:
		msg = "Команда не найдена"
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
}

func validateCommand(command []string, count int) bool {
	return len(command) == count
}

func toUpperSlice(commands []string) (words []string) {
	for _, w := range commands {
		words = append(words, strings.ToUpper(w))
	}
	return
}
