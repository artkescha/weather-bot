package main

import (
	"bots/telegram/weather_bot/resolver"
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"os"
)


func main() {

	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatalf("env PORT not exists")
		return
	}

	botToken, exists := os.LookupEnv("BOT_TOKEN")
	if !exists {
		log.Fatalf("env BOT_TOKEN not exists")
		return
	}

	weatherToken, exists := os.LookupEnv("WEATHER_TOKEN")
	if !exists {
		log.Fatalf("env WEATHER_TOKEN not exists")
		return
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("new bot failed: %s", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Debug = false
	bot.RemoveWebhook()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	baseUrl, exists := os.LookupEnv("BASE_URL")
	if !exists {
		log.Fatalf("env BASE_URL not exists")
		return
	}

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(baseUrl))

	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/")

	resolver := resolver.New(ctx, weatherToken, updates)
	resolver.Start(func(chatID int64, messageID int, message string) error {
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ReplyToMessageID = messageID
		_, err := bot.Send(msg)
		return err
	})

	http.ListenAndServe("0.0.0.0:"+port, nil)
}
