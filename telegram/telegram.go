package telegram

import (
	"log"
	"../config"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *telebot.BotAPI

func Connect() {
	var err error
	config := config.GetConfig()
	bot, err = telebot.NewBotAPI(config.TelegramApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)


}

func Poll() {
	u := telebot.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}

func SendToMe(s string) {
	log.Println("send message")
	msg := telebot.NewMessage(74450318, s)
	bot.Send(msg)
}

func Send(s string) {
	log.Println("send message")
	msg := telebot.NewMessage(-23576602, s)
	bot.Send(msg)
}
