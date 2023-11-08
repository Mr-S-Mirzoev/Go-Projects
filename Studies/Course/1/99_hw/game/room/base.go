package room

import (
	"fmt"

	"gitlab.com/sergeymirzoev/lectures-2021-2/1/99_hw/game/inventory"
)

// Implementation of RoomBase abstraction
type RoomBase struct {
	Connections []string
	Items       map[string]inventory.Item
	Name        string
	DoorClosed  bool
}

func (r *RoomBase) UseDoor() (string, error) {
	return "", fmt.Errorf("не к чему применить")
}

func (r *RoomBase) TakeItem(itemName string) (inventory.Item, error) {
	item, mItemExist := r.Items[itemName]
	if !mItemExist {
		return inventory.Item{}, fmt.Errorf("нет такого")
	}

	delete(r.Items, itemName)
	return item, nil
}

func (r *RoomBase) PutItem(item inventory.Item) {
	r.Items[item.Name] = item
}

func (r *RoomBase) RemoveTask(task string) error {
	return fmt.Errorf("NoTaskListError: не на кухне")
}

func (r *RoomBase) getName() string {
	return r.Name
}

func (r *RoomBase) HasItem(itemName string) bool {
	for i := range r.Items {
		if r.Items[i].Name == itemName {
			return true
		}
	}

	return false
}

func (r *RoomBase) HasPath(roomName string) (bool, error) {
	for i := range r.Connections {
		if r.Connections[i] == roomName {
			if roomName == "улица" && r.DoorClosed {
				return false, fmt.Errorf("дверь закрыта")
			} else {
				return true, nil
			}
		}
	}

	return false, fmt.Errorf("нет пути в %s", roomName)
}
