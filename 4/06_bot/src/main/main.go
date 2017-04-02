package main

// go get gopkg.in/telegram-bot-api.v4
// можем же собирать не у себя  и можем использовать у плохого разработчика
// скачиваем мы не у проверенных лиц

// способны решить проблему с пакетами
// gb
// go vendor
// название пакте - название исполняемго файла

// как создать зависимость
// go get -u github.com/constabulary/gb/..
// go get gopkg.in/telegram-bot-api.v4
// добавить в PATH - $GOPATH/bin
// export PATH=$PATH:$GOPATH/bin
// gb vendor fetch gopkg.in/telegram-bot-api.v4

// gb vendor restore - удалили папку vendor и оставили папку manifest000
// и установил версии как написанно в manifest
// go vendor update  - при обновлении версии в mainfest

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// для вендоринга используется GB
// сборка проекта осуществляется с помощью gb build
// установка зависимостей - gb vendor fetch gopkg.in/telegram-bot-api.v4
// установка зависимостей из манифеста - gb vendor restore

type Joke struct {
	ID   uint32 `json:"id"`
	Joke string `json:"joke"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value Joke   `json:"value"`
}

var buttons = []tgbotapi.KeyboardButton{
	{Text: "Get Joke"},
}

// При старте приложения, оно скажет телеграму ходить с обновлениями по этому URL
const WebhookURL = "https://startbot0.herokuapp.com/"

func getJoke() string {
	c := http.Client{}
	resp, err := c.Get("http://api.icndb.com/jokes/random?limitTo=[nerdy]")
	if err != nil {
		return "jokes API not responding"
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	joke := JokeResponse{}

	err = json.Unmarshal(body, &joke)
	if err != nil {
		return "Joke error"
	}

	return joke.Value.Joke
}

func main() {
	// Неroku не знает по какому порту будет прокидываться приложение
	// так heroku понимает куда именно обащаться за приложением
	// Heroku прокидывает порт для приложения в переменную окружения PORT
	port := os.Getenv("PORT")
	// библиотека создаст структуру в которую положит access(доступ) the HTTP API
	// сделает запрос и возьмет имя и положит к себе
	bot, err := tgbotapi.NewBotAPI("300203760:AAH9cwD3NpdcB6PVpxfzKVSfUd_uwzhSYZE")
	if err != nil {
		log.Fatal(err)
	}
	// ставим debug есть, чтобы все сообщения шли в логи
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Устанавливаем вебхук
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		log.Fatal(err)
	}

	// надо поставить webhook, чтобы сказать боту - ходи пожалуйтся к нам по этому Url
	// ListenForWebhook - делает тоже самое что и http.HandleFunc(pattern, handler)
	// просто скрывает добавления и парсинга ответа и добавление его в канал
	updates := bot.ListenForWebhook("/")
	// надо повесить бота на такой-то порт - в этом потоке надо обрабатывать update
	go http.ListenAndServe(":"+port, nil)

	// получаем все обновления из канала updates
	for update := range updates {
		var message tgbotapi.MessageConfig
		log.Println("received text: ", update.Message.Text)

		switch update.Message.Text {
		case "Cp":
			message = tgbotapi.NewMessage(update.Message.Chat.ID, getJoke())
		case "Get Joke":
			// Если пользователь нажал на кнопку, то придёт сообщение "Get Joke"
			message = tgbotapi.NewMessage(update.Message.Chat.ID, getJoke())
		default:
			message = tgbotapi.NewMessage(update.Message.Chat.ID, `Press "Get Joke" to receive joke`)
		}

		// В ответном сообщении просим показать клавиатуру
		// чтобы можно было просто показать кнопку, а не писать текст
		// передаем слаис слайсов
		message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)

		bot.Send(message)
	}
}
