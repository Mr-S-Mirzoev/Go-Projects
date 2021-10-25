package inventory

import "fmt"

type Inventory interface {
	Add(Item) (string, error)
	Get(string) (Item, error)
	Check(string) bool
}

type InventoryImpl struct {
	Items map[string]Item
}

func (inv *InventoryImpl) Add(item Item) (string, error) {
	var actionResult string

	if item.Name == "рюкзак" {
		actionResult = "вы надели: рюкзак"
	} else {
		if _, mItemExist := inv.Items["рюкзак"]; !mItemExist {
			return "", fmt.Errorf("некуда класть")
		}

		item.Location = "рюкзак"
		actionResult = fmt.Sprintf("предмет добавлен в инвентарь: %s", item.Name)
	}

	inv.Items[item.Name] = item
	return actionResult, nil
}

func (inv *InventoryImpl) Get(itemName string) (Item, error) {
	item, mItemExist := inv.Items[itemName]
	if !mItemExist {
		return Item{}, fmt.Errorf("нет такого")
	}

	item.Location = ""
	return item, nil
}

func (inv *InventoryImpl) Check(itemName string) bool {
	_, mItemExist := inv.Items[itemName]
	return mItemExist
}
