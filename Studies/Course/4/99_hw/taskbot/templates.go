package main

import "fmt"

const assignedToMe = "\nassignee: я"

type ReplyTemplate interface {
	String(int, UserData) string
}

type TasksReplyTemplate struct {
	task     string
	assignee UserData
	assigned bool
	commands string
}

type ManipulationCurrentReplyTemplate struct {
	task string
}

type ManipulationOldReplyTemplate struct {
	task string
}

func CreateTasksReplyTemplate(t Task) *TasksReplyTemplate {
	tmpl := &TasksReplyTemplate{
		task: fmt.Sprintf("%d. %s by @%s", t.TaskID(), t.TaskDescription(), t.CreatedBy().UserNick),
	}
	taskAssignee, err := t.AssignedTo()
	if err != nil {
		tmpl.assigned = false
		tmpl.commands = fmt.Sprintf("\n/assign_%d", t.TaskID())
	} else {
		tmpl.assigned = true
		tmpl.commands = fmt.Sprintf("\n/unassign_%d /resolve_%d", t.TaskID(), t.TaskID())
	}
	tmpl.assignee = taskAssignee

	return tmpl
}

func CreateManipulationCurrentReplyTemplate(taskDescription string) *ManipulationCurrentReplyTemplate {
	return &ManipulationCurrentReplyTemplate{
		task: taskDescription,
	}
}

func CreateManipulationOldReplyTemplate(taskDescription string) *ManipulationOldReplyTemplate {
	return &ManipulationOldReplyTemplate{
		task: taskDescription,
	}
}

func (tmpl *ManipulationCurrentReplyTemplate) String(typeOfHandler int, curUser UserData) string {
	messageText := ""
	switch typeOfHandler {
	case ASSIGN:
		messageText = fmt.Sprintf(
			"Задача \"%s\" назначена на вас",
			tmpl.task,
		)
	case UNASSIGN:
		messageText = "Принято"
	case RESOLVE:
		messageText = fmt.Sprintf(
			"Задача \"%s\" выполнена",
			tmpl.task,
		)
	}
	return messageText
}

func (tmpl *TasksReplyTemplate) String(typeOfHandler int, curUser UserData) string {
	result := tmpl.task
	if tmpl.assigned {
		if typeOfHandler == ALL || typeOfHandler == MY {
			if tmpl.assignee == curUser {
				if typeOfHandler != MY {
					result += assignedToMe
				}
				result += tmpl.commands
			} else {
				result += fmt.Sprintf("\nassignee: @%s", tmpl.assignee.UserNick)
			}
		}
	} else {
		result += tmpl.commands
	}
	return result
}

func (tmpl *ManipulationOldReplyTemplate) String(typeOfHandler int, otherUser UserData) string {
	messageText := ""
	switch typeOfHandler {
	case ASSIGN:
		messageText = fmt.Sprintf(
			"Задача \"%s\" назначена на @%s",
			tmpl.task,
			otherUser.UserNick,
		)
	case UNASSIGN:
		messageText = fmt.Sprintf(
			"Задача \"%s\" осталась без исполнителя",
			tmpl.task,
		)
	case RESOLVE:
		messageText = fmt.Sprintf(
			"Задача \"%s\" выполнена @%s",
			tmpl.task,
			otherUser.UserNick,
		)
	}
	return messageText
}
