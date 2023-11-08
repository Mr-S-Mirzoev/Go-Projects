package room

import "strings"

type Corridor struct {
	RoomBase
}

func (cr *Corridor) LookAround() string {
	return "ничего интересного. можно пройти - " + strings.Join(cr.Connections, ", ")
}

func (cr *Corridor) Greet() string {
	return "ничего интересного. можно пройти - " + strings.Join(cr.Connections, ", ")
}

func (cr *Corridor) UseDoor() (string, error) {
	if cr.DoorClosed {
		cr.DoorClosed = false
		return "дверь открыта", nil
	} else {
		cr.DoorClosed = true
		return "дверь закрыта", nil
	}
}
