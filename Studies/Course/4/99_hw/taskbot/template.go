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
