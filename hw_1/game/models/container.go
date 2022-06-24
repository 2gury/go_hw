package models

type Container struct {
	Name         string
	OutputName   string
	ActionOutput string
	Items        []*Item
}

func (ci *Container) AddItem(item *Item) {
	ci.Items = append(ci.Items, item)
}

func (ci *Container) RemoveItem(item *Item) *Item {
	idItem := -1
	var curItem *Item
	for id, it := range ci.Items {
		if it == item {
			idItem = id
			curItem = it
		}
	}
	if idItem == -1 {
		return nil
	}
	ci.Items = append(ci.Items[:idItem], ci.Items[idItem+1:]...)
	return curItem
}
