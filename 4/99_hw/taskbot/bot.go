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
	WebhookURL = "https://dbd9-83-242-55-252.ngrok.io"
	BotToken   = "2146592204:AAHcvCly8lQJ0NxbS3mdy1uLtH1LGfDOJ48"
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

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if ctx == nil {
			_, err = w.Write([]byte("Troubles with context"))
		} else {
			_, err = w.Write([]byte("all is working"))
		}
		log.Fatalf("http err: %v\n", err)
	})

	useNgrok := os.Getenv("USE_NGROK")
	if useNgrok == "" {
		WebhookURL = "http://127.0.0.1:8081"
		BotToken = "_golangcourse_test"
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
		BotHandle: bot,
	}

	// получаем все обновления из канала updates
	for update := range updates {
		log.Printf("upd: %#v\n", update)

		if update.Message == nil {
			log.Printf("Got null message")
			continue
		}

		hdlr.handleMessage(update.Message)
	}

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
