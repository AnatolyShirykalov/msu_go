package main

// gb build main
// ./bin/main
// heroku logs --tail

import (
	"fmt"
	"runtime"
	"time"
	// "github.com/gamer"
	"gopkg.in/telegram-bot-api.v4"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

// // для вендоринга используется GB
// // сборка проекта осуществляется с помощью gb build
// // установка зависимостей - gb vendor fetch gopkg.in/telegram-bot-api.v4
// // установка зависимостей из манифеста - gb vendor restore

// При старте приложения, оно скажет телеграму ходить с обновлениями по этому URL
const WebhookURL = "https://gamebot0.herokuapp.com/bot"

var Gamers map[string]*Player
var buttons = [][]tgbotapi.KeyboardButton{
	{
		{Text: "осмотреться"},
	},
	// 	{Text: "идти коридор"},
	// 	{Text: "надеть рюкзак"},
	// }, {
	// 	{Text: "осмотреться"},
	// 	{Text: "осмотреться"},
	// 	{Text: "осмотреться"},
	// 	{Text: "осмотреться"},
	// 	{Text: "осмотреться"},
	// },
}

func main() {
	// Неroku не знает по какому порту будет прокидываться приложение
	// так heroku понимает куда именно обащаться за приложением
	// Heroku прокидывает порт для приложения в переменную окружения PORT
	port := os.Getenv("PORT")
	lastOutput := make(map[string]string)
	Gamers = make(map[string]*Player)
	// библиотека создаст структуру в которую положит access(доступ) the HTTP API
	// сделает запрос и возьмет имя и положит к себе
	bot, err := tgbotapi.NewBotAPI("357629161:AAGYMC-PpkQDohL8fI5OVIT1QLSWtjIQWZk")
	if err != nil {
		log.Fatal(err)
	}
	// ставим debug есть, чтобы все сообщения шли в логи
	// bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Устанавливаем вебху	к
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		log.Fatal(err)
	}
	// mu := &sync.Mutex{}

	// надо поставить webhook, чlastтобы сказать боту - ходи пожалуйтся к нам по этому Url
	// ListenForWebhook - делает тоже самое что и http.HandleFunc(pattern, handler)
	// просто скрывает добавления и парсинга ответа и добавление его в канал
	updates := bot.ListenForWebhook("/bot")
	// надо повесить бота на такой-то порт - в этом потоке надо обрабатывать update
	go http.ListenAndServe(":"+port, nil)

	// получаем все обновления из канала updates
	for update := range updates {
		var message tgbotapi.MessageConfig
		answer := ""
		if update.Message.Text == "DELETE PLAYERS" && update.Message.Chat.UserName == "albina_artist" {
			gamers := ""
			flag := false
			for name := range Gamers {
				if flag {
					gamers += ", " + name
				} else {
					gamers += name
					flag = true
				}
			}
			answer = fmt.Sprintf("вышло из игры %d: %s", len(Gamers), gamers)
			DeletePlayers()
		}
		if _, ok := Gamers[update.Message.Chat.FirstName]; !ok {
			if update.Message.Text == "/start" {
				answer = update.Message.Chat.FirstName + ", начнем играть 0_0"
			}
			Gamers[update.Message.Chat.FirstName] = NewPlayer(update.Message.Chat.FirstName)
			mu := &sync.Mutex{}
			go func() {
				for msg := range Gamers[update.Message.Chat.FirstName].GetOutput() {
					mu.Lock()
					lastOutput[update.Message.Chat.FirstName] = msg
					log.Printf("Authorized on account %s", update.Message.Chat.FirstName+msg)
					mu.Unlock()
				}
			}()
			if len(Gamers) == 1 {
				InitGame()
			}
			AddPlayer(Gamers[update.Message.Chat.FirstName])
		} else {
			if update.Message.Text == "/start" {
				answer = update.Message.Chat.FirstName + ", ты уж в игре)"
			}
		}
		if answer == "" {
			lastOutput = map[string]string{}
			Gamers[update.Message.Chat.FirstName].HandleInput(update.Message.Text)
			time.Sleep(time.Microsecond)
			runtime.Gosched() // дадим считать ответ
			answer = lastOutput[update.Message.Chat.FirstName]
		}
		message = tgbotapi.NewMessage(update.Message.Chat.ID, answer)
		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons[0])

		bot.Send(message)
	}
}
func DeletePlayers() {
	for r := range Gamers {
		Gamers[r] = nil
	}
	G.Rooms = InitRoom()
}
