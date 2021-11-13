package task_manager

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

}
func (mgr *TaskManagerInMemory) GetTasksAssignedToUser(assigneeNick string) []Task {

}
func (mgr *TaskManagerInMemory) GetTasksCreatedByUser(creatorNick string) []Task {

}
func (mgr *TaskManagerInMemory) Assign(taskId int, assigneeNick string) (string, bool, error) {

}
func (mgr *TaskManagerInMemory) Unassign(taskId int, assigneeNick string) error {

}
func (mgr *TaskManagerInMemory) Resolve(taskId int, resolverNick string) (string, bool, error) {

}
