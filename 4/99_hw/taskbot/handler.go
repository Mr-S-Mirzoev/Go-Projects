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

const (
	noTasks         = `Нет задач`
	notMyTask       = `Задача не на вас`
	notAssigned     = `Задача не назначена`
	ALL         int = iota
	MY
	OWNER
)

func (hdlr *Handler) handleTasks(user UserData, typeOfHandler int) (string, error) {
	var tasks []Task
	switch typeOfHandler {
	case ALL:
		tasks = hdlr.Mngr.GetAllTasks()
	case MY:
		tasks = hdlr.Mngr.GetTasksAssignedToUser(user.UserNick)
	case OWNER:
		tasks = hdlr.Mngr.GetTasksCreatedByUser(user.UserNick)
	default:
		return "", UnknownHandlerTypeError{
			HandlerType: typeOfHandler,
		}
	}

	if len(tasks) == 0 {
		return noTasks, nil
	}

	returnString := ""
	for i, task := range tasks {
		if i != 0 {
			returnString += "\n"
		}

		returnString += fmt.Sprintf("%d. %s by @%s\n", task.TaskId(), task.TaskDescription(), task.CreatedBy())

		if task.Assigned() {
			if typeOfHandler == ALL {
				assignee, _ := task.AssignedTo()
				if assignee == user {
					returnString += fmt.Sprintf("assignee: я\n")
				} else {
					returnString += fmt.Sprintf("assignee: @%s\n", assignee)
				}
			}

			returnString += fmt.Sprintf("/unassign_%d /resolve_%d\n", task.TaskId(), task.TaskId())
		} else {
			returnString += fmt.Sprintf("/assign_%d\n", task.TaskId())
		}
	}

	return returnString, nil
}

func (hdlr *Handler) handleAssign(taskId int, user UserData) (map[ChatID]string, error) {
	oldAssignee, shouldNotify, err := hdlr.Mngr.Assign(taskId, user)
	if err != nil {
		_, ok := err.(NoSuchIDError)
		if ok {
			return map[ChatID]string{
				user.ID: fmt.Sprintf("Нет такой задачи с номером: %d", taskId),
			}, nil
		}

		return nil, err
	}

	taskDescription, _ := hdlr.Mngr.GetTaskDescriptionByID(taskId)

	returnMap := map[ChatID]string{
		user.ID: fmt.Sprintf("Задача \"%s\" назначена на вас", taskDescription),
	}

	if shouldNotify {
		returnMap[oldAssignee.ID] = fmt.Sprintf(
			"Задача \"%s\" назначена на @%s",
			taskDescription,
			user.UserNick,
		)
	}

	return returnMap, nil
}

func (hdlr *Handler) handleUnassign(taskId int, user UserData) (map[ChatID]string, error) {
	err := hdlr.Mngr.Unassign(taskId, user)
	if err != nil {
		switch err.(type) {
		case NoSuchIDError:
			return map[ChatID]string{
				user.ID: fmt.Sprintf("Нет такой задачи с номером: %d", taskId),
			}, nil
		case NotMyTaskError:
			return map[ChatID]string{
				user.ID: notMyTask,
			}, nil
		case NotAssignedError:
			return map[ChatID]string{
				user.ID: notAssigned,
			}, nil
		default:
			return nil, err
		}
	}
}

func (hdlr *Handler) handleResolve(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleNew(message *tgbotapi.Message) (string, error) {
	return "", nil
}

func (hdlr *Handler) handleMessage(message *tgbotapi.Message) (map[ChatID]string, error) {
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

	switch message.Text {
	case "/tasks":
		result, err := hdlr.handleTasks(userData, ALL)
		return map[ChatID]string{
			userData.ID: result,
		}, err
	case "/my":
		result, err := hdlr.handleTasks(userData, MY)
		return map[ChatID]string{
			userData.ID: result,
		}, err
	case "/owner":
		result, err := hdlr.handleTasks(userData, OWNER)
		return map[ChatID]string{
			userData.ID: result,
		}, err
	default:
		commandWithArgs := strings.Split(message.Text, " ")
		if len(commandWithArgs) == 1 {
			commandWithArgs = strings.Split(commandWithArgs[0], "_")
			if len(commandWithArgs) != 2 {
				return nil, fmt.Errorf("Wrong command: %s", message.Text)
			}

			id, err := strconv.Atoi(commandWithArgs[1])
			if err != nil {
				return nil, fmt.Errorf("Failed to convert ID to int in: %s (%v)", message.Text, err)
			}

			command := commandWithArgs[0]
			switch command {
			case "/assign":
				return hdlr.handleAssign(id, userData)
			case "/unassign":
				return hdlr.handleUnassign(id, userData)
			case "/resolve":
				hdlr.handleResolve(message)
			default:
				return nil, fmt.Errorf("Unknown command with id: %s", command)
			}
		}

		if commandWithArgs[0] != "/new" {
			return nil, fmt.Errorf("Unknown command with multiple words: %s", commandWithArgs[0])
		}

		result, err := hdlr.handleNew(message)
		return []string{
			result,
		}, err
	}
}
