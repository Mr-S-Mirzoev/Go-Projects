package main

import "fmt"

type NoSuchIDError struct {
	TaskID int
}

func (err NoSuchIDError) Error() string {
	return fmt.Sprintf("No such task id: %d", err.TaskID)
}

type NotMyTaskError struct {
	CorrectAssignee string
}

func (err NotMyTaskError) Error() string {
	return fmt.Sprintf("Not your task. Belongs to: %s", err.CorrectAssignee)
}

type NotAssignedError struct {
	TaskID int
}

func (err NotAssignedError) Error() string {
	return fmt.Sprintf("Task %d hasn't been assigned yet", err.TaskID)
}

type UnknownHandlerTypeError struct {
	HandlerType int
}

func (err UnknownHandlerTypeError) Error() string {
	return fmt.Sprintf("Unknown handler type: %d", err.HandlerType)
}
