package room

import "gitlab.com/sergeymirzoev/lectures-2021-2/1/99_hw/game/inventory"

type RoomIFace interface {
	LookAround() string
	HasItem(string) bool
	TakeItem(string) (inventory.Item, error)
	PutItem(inventory.Item)
	HasPath(string) (bool, error)
	UseDoor() (string, error)
	Greet() string
	RemoveTask(string) error
	getName() string
}
