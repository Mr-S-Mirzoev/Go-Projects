package main

// сюда писать код

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	WebhookURL = "http://127.0.0.1:8081"
	BotToken   = "_golangcourse_test"
)

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return fmt.Errorf("newBotAPI failed: %s", err)
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	if ctx == nil {
		log.Fatal("Just to disable go-lint error")
	}

	// Get the PORT from env vars cause it's implicitly set by Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	hdlr := Handler{
		Mngr: &TaskManagerInMemory{
			Tasks:  make(map[int]Task),
			LastID: 1,
		},
	}

	// получаем все обновления из канала updates
	for update := range updates {
		log.Printf("upd: %#v\n", update)

		replies, err := hdlr.handleMessage(update.Message)
		if err != nil {
			log.Fatalf(
				"Произошла ошибка для чата %d при сообщении %s: %v",
				update.Message.Chat.ID,
				update.Message.Text,
				err,
			)
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Произошла неизвестная ошибка",
			)
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatalf(
					"Failed to send message: %v : %v",
					msg,
					err,
				)
			}
			continue
		}

		for chatId, messageText := range replies {
			msg := tgbotapi.NewMessage(
				int64(chatId),
				messageText,
			)
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatalf(
					"Failed to send message: %v : %v",
					msg,
					err,
				)
			}
		}
		continue
	}

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
