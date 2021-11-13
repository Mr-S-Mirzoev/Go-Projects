package main

type TaskManager interface {
	CreateTask(string, UserData, UserData) Task
	GetAllTasks() []Task
	GetTasksAssignedToUser(string) []Task
	GetTasksCreatedByUser(string) []Task
	Assign(int, UserData) (UserData, bool, error)
	Unassign(int, UserData) error
	Resolve(int, UserData) (UserData, bool, error)
	GetTaskDescriptionByID(int) (string, error)
}

type TaskManagerInMemory struct {
	tasks  map[int]Task
	lastId int
}

func (mgr *TaskManagerInMemory) CreateTask(taskDescription string, cretorNick, assigneeNick UserData) Task {
	newTask := &TaskStruct{
		ID:          mgr.lastId,
		Description: taskDescription,
		Creator:     cretorNick,
		Assignee:    assigneeNick,
	}

	mgr.tasks[mgr.lastId] = newTask
	mgr.lastId++

	return newTask
}

func (mgr *TaskManagerInMemory) GetAllTasks() []Task {
	taskList := make([]Task, 0, len(mgr.tasks))
	for _, task := range mgr.tasks {
		taskList = append(taskList, task)
	}
	return taskList
}

func (mgr *TaskManagerInMemory) GetTasksAssignedToUser(assigneeNick string) []Task {
	taskList := make([]Task, 0, 10)
	for _, task := range mgr.tasks {
		assignee, err := task.AssignedTo()
		if err == nil && assignee.UserNick == assigneeNick {
			taskList = append(taskList, task)
		}
	}
	return taskList
}

func (mgr *TaskManagerInMemory) GetTasksCreatedByUser(creatorNick string) []Task {
	taskList := make([]Task, 0, 10)
	for _, task := range mgr.tasks {
		if task.CreatedBy().UserNick == creatorNick {
			taskList = append(taskList, task)
		}
	}
	return taskList
}

func (mgr *TaskManagerInMemory) Assign(taskId int, assignee UserData) (UserData, bool, error) {
	task, ok := mgr.tasks[taskId]
	if !ok {
		return UserData{}, false, NoSuchIDError{
			TaskID: taskId,
		}
	}

	oldAssignee, shouldNotify := task.AssignTo(assignee)
	return oldAssignee, shouldNotify, nil
}

func (mgr *TaskManagerInMemory) Unassign(taskId int, assignee UserData) error {
	task, ok := mgr.tasks[taskId]
	if !ok {
		return NoSuchIDError{
			TaskID: taskId,
		}
	}

	taskAssignee, err := task.AssignedTo()
	if err != nil {
		return err
	}

	if taskAssignee != assignee {
		return NotMyTaskError{
			CorrectAssignee: taskAssignee.UserNick,
		}
	}

	task.AssignTo(UserData{})
	return nil
}

func (mgr *TaskManagerInMemory) Resolve(taskId int, resolver UserData) (UserData, bool, error) {
	task, ok := mgr.tasks[taskId]
	if !ok {
		return UserData{}, false, NoSuchIDError{
			TaskID: taskId,
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
	delete(mgr.tasks, taskId)

	if creator == resolver {
		return UserData{}, false, nil
	}
	return creator, true, nil
}

func (mgr *TaskManagerInMemory) GetTaskDescriptionByID(taskId int) (string, error) {
	task, ok := mgr.tasks[taskId]
	if !ok {
		return "", NoSuchIDError{
			TaskID: taskId,
		}
	}

	return task.TaskDescription(), nil
}
