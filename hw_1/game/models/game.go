package models

import (
	"fmt"
)

type Game struct {
	CurrentPlayer *Player
	Rooms         map[string]*Room
	Actions       map[string]func(*Player)
	Missions      []*Mission
}

func NewGame(player *Player, rooms []*Room) *Game {
	mapRooms := map[string]*Room{}
	for _, room := range rooms {
		mapRooms[room.Name] = room
	}
	return &Game{
		CurrentPlayer: player,
		Rooms:         mapRooms,
		Actions:       map[string]func(*Player){},
	}
}

func (g *Game) AddRoomTransfer(first, second *Room) {
	first.Rooms = append(first.Rooms, second.Name)
	second.Rooms = append(second.Rooms, first.Name)
	g.Rooms[first.Name] = first
	g.Rooms[second.Name] = second
}

func (g *Game) AddGlobalTransfer(first, second *Room) {
	first.Rooms = append(first.Rooms, second.GlobalTransfer)
	second.Rooms = append(second.Rooms, first.Name)
	g.Rooms[first.Name] = first
	g.Rooms[second.GlobalTransfer] = second
	g.Rooms[second.Name] = second
}

func (g *Game) GetMission(missionPurpose string) *Mission {
	for _, mission := range g.Missions {
		if mission.Purpose == missionPurpose {
			return mission
		}
	}
	return nil
}

func (g *Game) GetMissions() string {
	outputMessage := ""
	for id, mission := range g.Missions {
		if !mission.IsCompleted && id == len(g.Missions)-1 {
			outputMessage += fmt.Sprintf(" %s", mission.Purpose)
		} else if !mission.IsCompleted {
			outputMessage += fmt.Sprintf(" %s и", mission.Purpose)
		}
	}
	return outputMessage
}

func (g *Game) GetRoom(roomName string) *Room {
	room, ok := g.Rooms[roomName]
	if !ok {
		return nil
	}
	return room
}

func (g *Game) GetItem(itemName string) *Item {
	for _, container := range g.CurrentPlayer.CurrentRoom.Containers {
		for _, item := range container.Items {
			if item.Name == itemName {
				container.RemoveItem(item)
				return item
			}
		}
	}
	return nil
}

func (g *Game) GetContainer(containerName string) *Container {
	for _, container := range g.CurrentPlayer.CurrentRoom.Containers {
		if container.Name == containerName {
			return container
		}
	}
	return nil
}

func (g *Game) GetRooms() string {
	outputMessage := ""
	for id, room := range g.CurrentPlayer.CurrentRoom.Rooms {
		if id == len(g.CurrentPlayer.CurrentRoom.Rooms)-1 {
			outputMessage += fmt.Sprintf(" %s", room)
		} else {
			outputMessage += fmt.Sprintf(" %s,", room)
		}
	}
	return outputMessage
}

func (g *Game) SetMissions(missions []string) {
	for _, mission := range missions {
		g.Missions = append(g.Missions, &Mission{
			Purpose: mission,
		})
	}
}

func (g *Game) LookAround() string {
	outputMessage := ""
	if g.CurrentPlayer.CurrentRoom.OutputMessage != "" {
		outputMessage = fmt.Sprintf("%s, ", g.CurrentPlayer.CurrentRoom.OutputMessage)
	}
	if g.CurrentPlayer.CurrentRoom.GetCountItems() != 0 {
		itemCounter := 0
		for idCon, container := range g.CurrentPlayer.CurrentRoom.Containers {
			if len(container.Items) != 0 {
				if idCon == 0 {
					outputMessage += fmt.Sprintf("на %s:", container.OutputName)
				} else {
					outputMessage += fmt.Sprintf(" на %s:", container.OutputName)
				}
			}
			for _, item := range container.Items {
				itemCounter++
				if len(container.Items) == 1 {
					outputMessage += fmt.Sprintf(" %s", item.Name)
				} else {
					if itemCounter == g.CurrentPlayer.CurrentRoom.GetCountItems() {
						outputMessage += fmt.Sprintf(" %s", item.Name)
					} else {
						outputMessage += fmt.Sprintf(" %s,", item.Name)
					}
				}
			}
		}

	} else {
		outputMessage += g.CurrentPlayer.CurrentRoom.EmptyMessage
	}
	if g.CurrentPlayer.CurrentRoom.ShowMission {
		mission := g.GetMissions()
		if mission != "" {
			outputMessage += fmt.Sprintf(", надо%s", mission)
		}
	}

	outputMessage += fmt.Sprintf(". можно пройти -%s", g.GetRooms())

	return outputMessage
}

func (g *Game) Go(roomName string) string {
	room := g.GetRoom(roomName)
	if room == nil {
		return fmt.Sprintf("нет пути в %s", roomName)
	}
	if !g.CurrentPlayer.CurrentRoom.IsRoomAvailable(room.Name) {
		return fmt.Sprintf("нет пути в %s", room.Name)
	}
	trigger, ok := g.Actions[fmt.Sprintf("идти %s", roomName)]
	if ok {
		trigger(g.CurrentPlayer)
	}
	if room.IsEnabled {
		g.CurrentPlayer.CurrentRoom = room
		return fmt.Sprintf("%s. можно пройти -%s", room.EntryMessage, g.GetRooms())
	}
	return fmt.Sprintf(room.EnableMessage)
}

func (g *Game) TakeItem(itemName string) string {
	if !g.CurrentPlayer.CanTakeItems {
		return "некуда класть"
	}
	item := g.GetItem(itemName)
	if item == nil {
		return "нет такого"
	}
	g.CurrentPlayer.TakeItem(item)
	return fmt.Sprintf("предмет добавлен в инвентарь: %s", item.Name)
}

func (g *Game) PutItem(itemName string) string {
	item := g.GetItem(itemName)
	if item == nil {
		return "нет такого"
	}
	trigger, ok := g.Actions[fmt.Sprintf("надеть %s", item.Name)]
	if ok {
		trigger(g.CurrentPlayer)
	}
	g.CurrentPlayer.TakeItem(item)
	return fmt.Sprintf("вы надели: %s", item.Name)
}

func (g *Game) Apply(itemSubject, itemObject string) string {
	itemSub := g.CurrentPlayer.GetItem(itemSubject)
	if itemSub == nil {
		return fmt.Sprintf("нет предмета в инвентаре - %s", itemSubject)
	}
	itemObj := g.GetContainer(itemObject)
	if itemObj == nil {
		return "не к чему применить"
	}
	trigger, ok := g.Actions[fmt.Sprintf("%s %s", itemSubject, itemObject)]
	if ok {
		trigger(g.CurrentPlayer)
	}
	return fmt.Sprintf(itemObj.ActionOutput)
}
