package room

import (
	"sort"
	"strings"
)

// Implementation of LivingRoom
type LivingRoom struct {
	RoomBase
}

func (lr *LivingRoom) LookAround() string {
	var s string
	if len(lr.Items) == 0 {
		s += "пустая комната"
	} else {
		whatWhere := make(map[string][]string, 5)
		for itemName, item := range lr.Items {
			where := item.Location
			whatWhere[where] = append(whatWhere[where], itemName)
		}

		for location, items := range whatWhere {
			s += location + ": "
			sort.Strings(items)
			s += strings.Join(items, ", ") + ", "
		}

		s = s[:len(s)-2]
	}

	s += ". можно пройти - " + strings.Join(lr.Connections, ", ")
	return s
}

func (lr *LivingRoom) Greet() string {
	return "ты в своей комнате. можно пройти - " + strings.Join(lr.Connections, ", ")
}
