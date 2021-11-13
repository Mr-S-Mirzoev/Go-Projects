package main

type TaskManager interface {
	CreateTask(string, UserData) Task
	GetAllTasks() []Task
	GetTasksAssignedToUser(string) []Task
	GetTasksCreatedByUser(string) []Task
	Assign(int, UserData) (UserData, bool, error)
	Unassign(int, UserData) (UserData, bool, error)
	Resolve(int, UserData) (UserData, bool, error)
	GetTaskDescriptionByID(int) (string, error)
}

type TaskManagerInMemory struct {
	Tasks  map[int]Task
	LastID int
}

func (mgr *TaskManagerInMemory) CreateTask(taskDescription string, creatorNick UserData) Task {
	newTask := &TaskStruct{
		ID:          mgr.LastID,
		Description: taskDescription,
		Creator:     creatorNick,
	}

	mgr.Tasks[mgr.LastID] = newTask
	mgr.LastID++

	return newTask
}

func (mgr *TaskManagerInMemory) GetAllTasks() []Task {
	taskList := make([]Task, 0, len(mgr.Tasks))
	for _, task := range mgr.Tasks {
		taskList = append(taskList, task)
	}
	return taskList
}

func (mgr *TaskManagerInMemory) GetTasksAssignedToUser(assigneeNick string) []Task {
	taskList := make([]Task, 0, 10)
	for _, task := range mgr.Tasks {
		assignee, err := task.AssignedTo()
		if err == nil && assignee.UserNick == assigneeNick {
			taskList = append(taskList, task)
		}
	}
	return taskList
}

func (mgr *TaskManagerInMemory) GetTasksCreatedByUser(creatorNick string) []Task {
	taskList := make([]Task, 0, 10)
	for _, task := range mgr.Tasks {
		if task.CreatedBy().UserNick == creatorNick {
			taskList = append(taskList, task)
		}
	}
	return taskList
}

func (mgr *TaskManagerInMemory) Assign(taskID int, assignee UserData) (UserData, bool, error) {
	task, ok := mgr.Tasks[taskID]
	if !ok {
		return UserData{}, false, NoSuchIDError{
			TaskID: taskID,
		}
	}

	oldAssignee, shouldNotify := task.AssignTo(assignee)
	return oldAssignee, shouldNotify, nil
}

func (mgr *TaskManagerInMemory) Unassign(taskID int, assignee UserData) (UserData, bool, error) {
	task, ok := mgr.Tasks[taskID]
	if !ok {
		return UserData{}, false, NoSuchIDError{
			TaskID: taskID,
		}
	}

	taskAssignee, err := task.AssignedTo()
	if err != nil {
		return UserData{}, false, err
	}

	if taskAssignee != assignee {
		return UserData{}, false, NotMyTaskError{
			CorrectAssignee: taskAssignee.UserNick,
		}
	}

	task.AssignTo(UserData{})
	creator := task.CreatedBy()
	if creator != assignee {
		return creator, true, nil
	}

	return UserData{}, false, nil
}

func (mgr *TaskManagerInMemory) Resolve(taskID int, resolver UserData) (UserData, bool, error) {
	task, ok := mgr.Tasks[taskID]
	if !ok {
		return UserData{}, false, NoSuchIDError{
			TaskID: taskID,
		}
	}

	taskAssignee, err := task.AssignedTo()
	if err != nil {
		return UserData{}, false, err
	}

	if taskAssignee != resolver {
		return UserData{}, false, NotMyTaskError{
			CorrectAssignee: taskAssignee.UserNick,
		}
	}

	creator := task.CreatedBy()
	delete(mgr.Tasks, taskID)

	if creator == resolver {
		return UserData{}, false, nil
	}
	return creator, true, nil
}

func (mgr *TaskManagerInMemory) GetTaskDescriptionByID(taskID int) (string, error) {
	task, ok := mgr.Tasks[taskID]
	if !ok {
		return "", NoSuchIDError{
			TaskID: taskID,
		}
	}

	return task.TaskDescription(), nil
}
