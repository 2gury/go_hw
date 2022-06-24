package models

type Room struct {
	Name           string
	OutputMessage  string
	EmptyMessage   string
	EntryMessage   string
	EnableMessage  string
	GlobalTransfer string
	ShowMission    bool
	Containers     []*Container
	Rooms          []string
	IsEnabled      bool
}

func (r *Room) IsRoomAvailable(roomName string) bool {
	for _, room := range r.Rooms {
		if room == roomName {
			return true
		}
	}
	return false
}

func (r *Room) AddContainer(container *Container) {
	r.Containers = append(r.Containers, container)
}

func (r *Room) GetContainer(containerName string) *Container {
	for _, container := range r.Containers {
		if container.Name == containerName {
			return container
		}
	}
	return nil
}

func (r *Room) GetCountItems() int {
	countItems := 0
	for _, container := range r.Containers {
		countItems += len(container.Items)
	}
	return countItems
}
