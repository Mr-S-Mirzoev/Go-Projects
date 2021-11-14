package main

import (
	"fmt"
	"sort"
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

		returnString += fmt.Sprintf("%d. %s by @%s", task.TaskID(), task.TaskDescription(), task.CreatedBy().UserNick)

		if task.Assigned() {
			if typeOfHandler == ALL || typeOfHandler == MY {
				assignee, err := task.AssignedTo()
				if err != nil {
					fmt.Printf(
						"Somehow error with assigned to occured even though taskId: %d was checked: %v",
						task.TaskID(),
						err,
					)
				}

				if assignee == user {
					if typeOfHandler != MY {
						returnString += "\nassignee: я"
					}
					returnString += fmt.Sprintf("\n/unassign_%d /resolve_%d", task.TaskID(), task.TaskID())
				} else {
					returnString += fmt.Sprintf("\nassignee: @%s", assignee.UserNick)
				}
			}
		} else {
			returnString += fmt.Sprintf("\n/assign_%d", task.TaskID())
		}
	}

	return returnString, nil
}

func (hdlr *Handler) handleAssign(taskID int, user UserData) (map[ChatID]string, error) {
	oldAssignee, shouldNotify, err := hdlr.Mngr.Assign(taskID, user)
	if err != nil {
		_, ok := err.(NoSuchIDError)
		if ok {
			return map[ChatID]string{
				user.ID: fmt.Sprintf("Нет такой задачи с номером: %d", taskID),
			}, nil
		}

		return nil, err
	}

	taskDescription, err := hdlr.Mngr.GetTaskDescriptionByID(taskID)
	if err != nil {
		fmt.Printf(
			"Somehow error with get description occured even though taskId: %d was checked: %v",
			taskID,
			err,
		)
	}

	returnMap := map[ChatID]string{
		user.ID: fmt.Sprintf("Задача \"%s\" назначена на вас", taskDescription),
	}

	if shouldNotify && user.ID != oldAssignee.ID {
		returnMap[oldAssignee.ID] = fmt.Sprintf(
			"Задача \"%s\" назначена на @%s",
			taskDescription,
			user.UserNick,
		)
	}

	return returnMap, nil
}

// TODO: Merge Unassign and Resolve handlers
func (hdlr *Handler) handleUnassign(taskID int, user UserData) (map[ChatID]string, error) {
	creator, shouldNotify, err := hdlr.Mngr.Unassign(taskID, user)
	if err != nil {
		switch err.(type) {
		case NoSuchIDError:
			return map[ChatID]string{
				user.ID: fmt.Sprintf("Нет такой задачи с номером: %d", taskID),
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

	returnMap := map[ChatID]string{user.ID: "Принято"}
	if shouldNotify {
		taskDescription, err := hdlr.Mngr.GetTaskDescriptionByID(taskID)
		if err != nil {
			fmt.Printf(
				"Somehow error with get description occured even though taskID: %d was checked: %v",
				taskID,
				err,
			)
		}
		returnMap[creator.ID] = fmt.Sprintf("Задача \"%s\" осталась без исполнителя", taskDescription)
	}

	return returnMap, nil
}

func (hdlr *Handler) handleResolve(taskID int, user UserData) (map[ChatID]string, error) {
	taskDescription, err := hdlr.Mngr.GetTaskDescriptionByID(taskID)
	if err != nil {
		fmt.Printf(
			"Error with get description for taskID %d occured, will be handled afterwards: %v",
			taskID,
			err,
		)
	}
	creator, shouldNotify, err := hdlr.Mngr.Resolve(taskID, user)
	if err != nil {
		switch err.(type) {
		case NoSuchIDError:
			return map[ChatID]string{
				user.ID: fmt.Sprintf("Нет такой задачи с номером: %d", taskID),
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

	returnMap := map[ChatID]string{
		user.ID: fmt.Sprintf("Задача \"%s\" выполнена", taskDescription),
	}
	if shouldNotify {
		returnMap[creator.ID] = fmt.Sprintf("Задача \"%s\" выполнена @%s", taskDescription, user.UserNick)
	}

	return returnMap, nil
}

func (hdlr *Handler) handleNew(taskDescription string, user UserData) string {
	task := hdlr.Mngr.CreateTask(taskDescription, user)
	return fmt.Sprintf("Задача \"%s\" создана, id=%d", task.TaskDescription(), task.TaskID())
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
				return nil, fmt.Errorf("wrong command: %s", message.Text)
			}

			id, err := strconv.Atoi(commandWithArgs[1])
			if err != nil {
				return nil, fmt.Errorf("failed to convert ID to int in: %s (%v)", message.Text, err)
			}

			command := commandWithArgs[0]
			switch command {
			case "/assign":
				return hdlr.handleAssign(id, userData)
			case "/unassign":
				return hdlr.handleUnassign(id, userData)
			case "/resolve":
				return hdlr.handleResolve(id, userData)
			default:
				return nil, fmt.Errorf("unknown command with id: %s", command)
			}
		}

		if commandWithArgs[0] != "/new" {
			return nil, fmt.Errorf("unknown command with multiple words: %s", commandWithArgs[0])
		}

		result := hdlr.handleNew(strings.Join(commandWithArgs[1:], " "), userData)
		return map[ChatID]string{
			userData.ID: result,
		}, nil
	}
}
