package main

import tgbotapi "gopkg.in/telegram-bot-api.v4"

type ChatID int

type UserData struct {
	ID       ChatID
	UserNick string
}

func FromTelegramMessage(message *tgbotapi.Message) UserData {
	return UserData{
		ID:       message.From.ID,
		UserNick: message.From.UserName,
	}
}
