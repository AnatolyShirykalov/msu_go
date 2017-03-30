package main

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI("300203760:AAH9cwD3NpdcB6PVpxfzKVSfUd_uwzhSYZE")
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		log.Fatal(err)
	}
	allPlayers := make(map[string]bool)
	if !allPlayers[bot.Self.UserName] {
		allPlayers[bot.Self.UserName] = true
	}
	players := map[string]*Player{
		bot: NewPlayer("Tristan"),
	}
	// if (G.Rooms[]bot)
}
