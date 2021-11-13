package main

type Task interface {
	// Returns if we have to notify anyone or not
	AssignTo(UserData) (UserData, bool)
	Assigned() bool
	AssignedTo() (UserData, error)
	CreatedBy() UserData
	TaskDescription() string
	TaskId() int
}

type TaskStruct struct {
	ID          int
	Description string
	Creator     UserData
	Assignee    UserData
}

func (t *TaskStruct) AssignTo(newAssignee UserData) (UserData, bool) {
	oldAssignee := t.Assignee
	if !t.Assigned() || newAssignee == oldAssignee {
		t.Assignee = newAssignee
		return UserData{}, false
	}

	t.Assignee = newAssignee
	return oldAssignee, true
}

func (t *TaskStruct) Assigned() bool {
	return t.Assignee.UserNick != ""
}

func (t *TaskStruct) AssignedTo() (UserData, error) {
	if !t.Assigned() {
		return UserData{}, NotAssignedError{
			TaskID: t.ID,
		}
	}

	return t.Assignee, nil
}

func (t *TaskStruct) CreatedBy() UserData {
	return t.Creator
}

func (t *TaskStruct) TaskDescription() string {
	return t.Description
}

func (t *TaskStruct) TaskId() int {
	return t.ID
}
