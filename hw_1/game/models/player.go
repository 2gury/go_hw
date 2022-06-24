package models

type Player struct {
	CurrentRoom  *Room
	Items        []*Item
	CanTakeItems bool
}

func (p *Player) TakeItem(item *Item) {
	p.Items = append(p.Items, item)
}

func (p *Player) GetItem(itemName string) *Item {
	for _, item := range p.Items {
		if item.Name == itemName {
			return item
		}
	}
	return nil
}
