package main

// сюда писать код

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tmngr "gitlab.com/sergeymirzoev/lectures-2021-2/4/99_hw/taskbot/task_manager"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "XXX"

	// урл выдаст вам нгрок или хероку
	WebhookURL = "https://525f2cb5.ngrok.io"
)

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return fmt.Errorf("NewBotAPI failed: %s", err)
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	// Get the PORT from env vars cause it's implicitly set by Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	hdlr := Handler{
		Mngr: tmngr.TaskManager,
	}

	// получаем все обновления из канала updates
	for update := range updates {
		log.Printf("upd: %#v\n", update)

		reply, err := handleMessage(update.Message)
		if err != nil {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				`there is only Habr feed availible`,
			)

			msg.ReplyMarkup = &tgbotapi.ReplyKeyboardMarkup{
				Keyboard: [][]tgbotapi.KeyboardButton{
					{
						{
							Text: "Habr",
						},
					},
				},
			}
			bot.Send(msg)
			continue
		}
	}

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
