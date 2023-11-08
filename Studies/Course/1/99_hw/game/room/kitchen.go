package room

import (
	"strings"
)

type Kitchen struct {
	RoomBase
	Tasks []string
}

func (k *Kitchen) LookAround() string {
	var s string
	if len(k.Tasks)+len(k.Items) == 0 {
		s += "ничего интересного"
	} else {
		if len(k.Items) == 0 {
			s += "пустая кухня"
		} else {
			s += "ты находишься на кухне, "
			whatWhere := make(map[string][]string, 5)
			for itemName, item := range k.Items {
				where := item.Location
				whatWhere[where] = append(whatWhere[where], itemName)
			}

			for location, items := range whatWhere {
				s += location + ": "
				s += strings.Join(items, ", ") + ", "
			}

			s = s[:len(s)-2]
		}

		if len(k.Tasks) == 0 {
			s += "."
		} else {
			s += ", надо "
			s += strings.Join(k.Tasks, " и ")
		}
	}

	s += ". можно пройти - " + strings.Join(k.Connections, ", ")
	return s
}

func (k *Kitchen) Greet() string {
	return "кухня, ничего интересного. можно пройти - " + strings.Join(k.Connections, ", ")
}

func (k *Kitchen) RemoveTask(task string) error {
	for i := range k.Tasks {
		if k.Tasks[i] == task {
			copy(k.Tasks[i:], k.Tasks[i+1:])
			k.Tasks[len(k.Tasks)-1] = "" // or the zero value of T
			k.Tasks = k.Tasks[:len(k.Tasks)-1]

			return nil
		}
	}

	return nil
}
