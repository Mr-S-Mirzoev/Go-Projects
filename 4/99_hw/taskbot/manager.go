package main

type TaskManager interface {
	CreateTask(string, string, string) Task
	GetAllTasks() []Task
	GetTasksAssignedToUser(string) []Task
	GetTasksCreatedByUser(string) []Task
	Assign(int, string) (string, bool, error)
	Unassign(int, string) error
	Resolve(int, string) (string, bool, error)
}

type TaskManagerInMemory struct {
	tasks  map[int]Task
	lastId int
}

func (mgr *TaskManagerInMemory) CreateTask(taskDescription, cretorNick, assigneeNick string) Task {
	ts := TaskStruct{
		ID:          mgr.lastId,
		Description: taskDescription,
		Creator:     cretorNick,
		Assignee:    assigneeNick,
	}
	mgr.tasks[mgr.lastId] = ts
}

func (mgr *TaskManagerInMemory) GetAllTasks() []Task {
	return nil
}

func (mgr *TaskManagerInMemory) GetTasksAssignedToUser(assigneeNick string) []Task {
	return nil
}

func (mgr *TaskManagerInMemory) GetTasksCreatedByUser(creatorNick string) []Task {
	return nil
}

func (mgr *TaskManagerInMemory) Assign(taskId int, assigneeNick string) (string, bool, error) {
	return "", false, nil
}

func (mgr *TaskManagerInMemory) Unassign(taskId int, assigneeNick string) error {
	return nil
}

func (mgr *TaskManagerInMemory) Resolve(taskId int, resolverNick string) (string, bool, error) {
	return "", false, nil
}
