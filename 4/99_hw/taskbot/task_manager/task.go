package taskmanager

import "fmt"

type Task interface {
	// Returns if we have to notify anyone or not
	AssignTo(string) (string, bool)
	Assigned() bool
	AssignedTo() (string, error)
	CreatedBy() string
	TaskDescription() string
	TaskId() int
}

type TaskStruct struct {
	ID          int
	Description string
	Creator     string
	Assignee    string
}

func (t *TaskStruct) AssignTo(newAssignee string) (string, bool) {
	oldAssignee := t.Assignee
	if !t.Assigned() || newAssignee == oldAssignee {
		t.Assignee = newAssignee
		return "", false
	}

	t.Assignee = newAssignee
	return oldAssignee, true
}

func (t *TaskStruct) Assigned() bool {
	return t.Assignee != ""
}

func (t *TaskStruct) AssignedTo() (string, error) {
	if !t.Assigned() {
		return "", fmt.Errorf("Task %d hasn't been assigned yet", t.ID)
	}

	return t.Assignee, nil
}

func (t *TaskStruct) CreatedBy() string {
	return t.Creator
}

func (t *TaskStruct) TaskDescription() string {
	return t.Description
}

func (t *TaskStruct) TaskId() int {
	return t.ID
}
