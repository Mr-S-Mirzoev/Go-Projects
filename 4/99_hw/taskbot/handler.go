package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Handler struct {
	Mngr      TaskManager
	BotHandle *tgbotapi.BotAPI
}

const (
	noTasks     = `Нет задач`
	notMyTask   = `Задача не на вас`
	notAssigned = `Задача не назначена`
)

const (
	ALL int = iota
	MY
	OWNER
)

const (
	ASSIGN int = iota
	UNASSIGN
	RESOLVE
)

func commandToHandlerType(command string) int {
	switch command {
	case "/assign":
		return ASSIGN
	case "/tasks":
		return ALL
	case "/my":
		return MY
	case "/unassign":
		return UNASSIGN
	case "/owner":
		return OWNER
	case "/resolve":
		return RESOLVE
	}
	return -1 // unreached code
}

func actOnError(err error, taskID int) (string, error) {
	switch err.(type) {
	case NoSuchIDError:
		return fmt.Sprintf("Нет такой задачи с номером: %d", taskID), nil
	case NotMyTaskError:
		return notMyTask, nil
	case NotAssignedError:
		return notAssigned, nil
	default:
		return "", err
	}
}

func (hdlr *Handler) sendMessage(chatID ChatID, messageText string) {
	msg := tgbotapi.NewMessage(
		int64(chatID),
		messageText,
	)

	_, err := hdlr.BotHandle.Send(msg)
	if err != nil {
		log.Fatalf(
			"Failed to send message: %v : %v",
			msg,
			err,
		)
	}
}

func (hdlr *Handler) handleTasks(user UserData, typeOfGetter int) error {
	var tasks []Task

	switch typeOfGetter {
	case ALL:
		tasks = hdlr.Mngr.GetAllTasks()
	case MY:
		tasks = hdlr.Mngr.GetTasksAssignedToUser(user.UserNick)
	case OWNER:
		tasks = hdlr.Mngr.GetTasksCreatedByUser(user.UserNick)
	default:
		return UnknownHandlerTypeError{
			HandlerType: typeOfGetter,
		}
	}

	if len(tasks) == 0 {
		hdlr.sendMessage(user.ID, noTasks)
		return nil
	}

	returnString := ""
	sortedTasks := make([]Task, 0, len(tasks))
	sortedTasks = append(sortedTasks, tasks...)
	sort.Slice(sortedTasks, func(i, j int) bool { return sortedTasks[i].TaskID() < sortedTasks[j].TaskID() })

	first := true
	for _, task := range sortedTasks {
		if !first {
			returnString += "\n\n"
		} else {
			first = false
		}

		replyTemplate := CreateTasksReplyTemplate(task)
		returnString += replyTemplate.String(typeOfGetter, user)
	}

	hdlr.sendMessage(user.ID, returnString)
	return nil
}

func (hdlr *Handler) handleTaskManipulation(taskID int, user UserData, typeOfTaskManipulation int) error {
	var err error
	var oldAssignee UserData
	var shouldNotify bool

	taskDescription, err := hdlr.Mngr.GetTaskDescriptionByID(taskID)
	if err != nil {
		fmt.Printf("Imposible error occured: %v", err)
	}

	switch typeOfTaskManipulation {
	case ASSIGN:
		oldAssignee, shouldNotify, err = hdlr.Mngr.Assign(taskID, user)
	case UNASSIGN:
		oldAssignee, shouldNotify, err = hdlr.Mngr.Unassign(taskID, user)
	case RESOLVE:
		oldAssignee, shouldNotify, err = hdlr.Mngr.Resolve(taskID, user)
	default:
		return UnknownHandlerTypeError{
			HandlerType: typeOfTaskManipulation,
		}
	}

	var messageText string
	if err != nil {
		messageText, err = actOnError(err, taskID)
		if err != nil {
			return err
		}

		hdlr.sendMessage(user.ID, messageText)
		return nil
	}

	switch typeOfTaskManipulation {
	case ASSIGN:
		messageText = fmt.Sprintf(
			"Задача \"%s\" назначена на вас",
			taskDescription,
		)
	case UNASSIGN:
		messageText = "Принято"
	case RESOLVE:
		messageText = fmt.Sprintf(
			"Задача \"%s\" выполнена",
			taskDescription,
		)
	default:
		return UnknownHandlerTypeError{
			HandlerType: typeOfTaskManipulation,
		}
	}

	hdlr.sendMessage(user.ID, messageText)

	if shouldNotify && user.ID != oldAssignee.ID {
		switch typeOfTaskManipulation {
		case ASSIGN:
			messageText = fmt.Sprintf(
				"Задача \"%s\" назначена на @%s",
				taskDescription,
				user.UserNick,
			)
		case UNASSIGN:
			messageText = fmt.Sprintf(
				"Задача \"%s\" осталась без исполнителя",
				taskDescription,
			)
		case RESOLVE:
			messageText = fmt.Sprintf(
				"Задача \"%s\" выполнена @%s",
				taskDescription,
				user.UserNick,
			)
		default:
			return UnknownHandlerTypeError{
				HandlerType: typeOfTaskManipulation,
			}
		}

		hdlr.sendMessage(oldAssignee.ID, messageText)
	}

	return nil
}

func (hdlr *Handler) handleNew(taskDescription string, user UserData) {
	task := hdlr.Mngr.CreateTask(taskDescription, user)
	messageText := fmt.Sprintf("Задача \"%s\" создана, id=%d", task.TaskDescription(), task.TaskID())
	hdlr.sendMessage(user.ID, messageText)
}

func (hdlr *Handler) handleUnknownError(err error, chatID ChatID, messageText string) {
	log.Printf(
		"Произошла ошибка для чата %d при сообщении %s: %v",
		chatID,
		messageText,
		err,
	)
	hdlr.sendMessage(chatID, "Произошла неизвестная ошибка")
}

func (hdlr *Handler) handleMessage(message *tgbotapi.Message) {
	/*
	* `/tasks`
	* `/new XXX YYY ZZZ` - создаёт новую задачу
	* `/assign_$ID` - делает пользователя исполнителем задачи
	* `/unassign_$ID` - снимает задачу с текущего исполнителя
	* `/resolve_$ID` - выполняет задачу, удаляет её из списка
	* `/my` - показывает задачи, которые назначены на меня
	* `/owner` - показывает задачи которые были созданы мной
	 */
	userData := FromTelegramMessage(message)

	if message.Text == "/tasks" || message.Text == "/my" || message.Text == "/owner" {
		err := hdlr.handleTasks(
			userData,
			commandToHandlerType(message.Text),
		)

		if err != nil {
			hdlr.handleUnknownError(
				err,
				ChatID(message.Chat.ID),
				message.Text,
			)
		}
	} else {
		if !strings.Contains(message.Text, " ") {
			commandWithArgs := strings.Split(message.Text, "_")

			if len(commandWithArgs) != 2 {
				hdlr.handleUnknownError(
					fmt.Errorf("wrong command: %s", message.Text),
					ChatID(message.Chat.ID),
					message.Text,
				)
				return
			}

			id, err := strconv.Atoi(commandWithArgs[1])
			if err != nil {
				hdlr.handleUnknownError(
					fmt.Errorf("failed to convert ID to int in: %s (%v)", message.Text, err),
					ChatID(message.Chat.ID),
					message.Text,
				)
				return
			}

			command := commandWithArgs[0]
			if command == "/assign" || command == "/unassign" || command == "/resolve" {
				err = hdlr.handleTaskManipulation(id, userData, commandToHandlerType(command))
			} else {
				err = fmt.Errorf("unknown command with id: %s", command)
			}

			if err != nil {
				hdlr.handleUnknownError(
					err,
					ChatID(message.Chat.ID),
					message.Text,
				)
				return
			}
		} else {
			commandWithArgs := strings.Split(message.Text, " ")

			if commandWithArgs[0] != "/new" {
				hdlr.handleUnknownError(
					fmt.Errorf("unknown command with multiple words: %s", commandWithArgs[0]),
					ChatID(message.Chat.ID),
					message.Text,
				)
			} else {
				hdlr.handleNew(strings.Join(commandWithArgs[1:], " "), userData)
			}
		}
	}
}
