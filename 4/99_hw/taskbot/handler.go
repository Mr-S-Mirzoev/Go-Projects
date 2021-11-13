package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Handler struct {
	Mngr TaskManager
}

func (hdlr *Handler) handleTasks(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleMy(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleOwner(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleAssign(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleUnassign(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleResolve(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleNew(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleMessage(message *tgbotapi.Message) (string, error) {
	/*
	* `/tasks`
	* `/new XXX YYY ZZZ` - создаёт новую задачу
	* `/assign_$ID` - делает пользователя исполнителем задачи
	* `/unassign_$ID` - снимает задачу с текущего исполнителя
	* `/resolve_$ID` - выполняет задачу, удаляет её из списка
	* `/my` - показывает задачи, которые назначены на меня
	* `/owner` - показывает задачи которые были созданы мной
	 */
	userNick := message.From.UserName

	switch message.Text {
	case "/tasks":
		return handleTasks(message)
	case "/my":
		return handleMy(message)
	case "/owner":
		return handleOwner(message)
	default:
		commandWithArgs := strings.Split(message.Text, " ")
		if len(commandWithArgs) == 1 {
			commandWithArgs = strings.Split(commandWithArgs[0], "_")
			if len(commandWithArgs) != 2 {
				return "", fmt.Errorf("Wrong command: %s", message.Text)
			}

			id, err := strconv.Atoi(commandWithArgs[1])
			if err != nil {
				return "", fmt.Errorf("Failed to convert ID to int in: %s (%v)", message.Text, err)
			}

			command := commandWithArgs[0]
			switch command {
			case "/assign":
				handleAssign(message)
			case "/unassign":
				handleUnassign(message)
			case "/resolve":
				handleResolve(message)
			default:
				return "", fmt.Errorf("Unknown command with id: %s", command)
			}
		}

		if commandWithArgs[0] != "/new" {
			return "", fmt.Errorf("Unknown command with multiple words: %s", commandWithArgs[0])
		}

		return handleNew(message)
	}
}
