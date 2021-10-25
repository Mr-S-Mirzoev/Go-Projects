package inventory

import "fmt"

type Item struct {
	Location     string
	Name         string
	Applications []string
}

func (it *Item) Use(onWhat string) error {
	for i := range it.Applications {
		if it.Applications[i] == onWhat {
			return nil
		}
	}

	return fmt.Errorf("UnapplicableError: item \"%s\" cannot be used on \"%s\"", it.Name, onWhat)
}
