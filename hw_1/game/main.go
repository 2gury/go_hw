package main

import (
	"gitlab.com/mailru-go/lectures-2022-1/01_intro/99_hw/game/models"
	"strings"
)

var CurrentGame *models.Game

func main() {
	initGame()
}

func initGame() {
	tableKitchen := &models.Container{
		Name:       "стол",
		OutputName: "столе",
	}
	tableRoom := &models.Container{
		Name:       "стол",
		OutputName: "столе",
	}
	door := &models.Container{
		Name:         "дверь",
		OutputName:   "двери",
		ActionOutput: "дверь открыта",
	}
	chair := &models.Container{
		Name:       "стул",
		OutputName: "стуле",
	}
	keys := &models.Item{Name: "ключи"}
	tea := &models.Item{Name: "чай"}
	backpack := &models.Item{Name: "рюкзак"}
	conspects := &models.Item{Name: "конспекты"}
	tableKitchen.AddItem(tea)
	tableRoom.AddItem(keys)
	tableRoom.AddItem(conspects)
	chair.AddItem(backpack)
	kitchen := &models.Room{
		Name:          "кухня",
		OutputMessage: "ты находишься на кухне",
		EmptyMessage:  "пустая кухня",
		EntryMessage:  "кухня, ничего интересного",
		ShowMission:   true,
		Containers:    []*models.Container{tableKitchen},
		IsEnabled:     true,
	}
	room := &models.Room{
		Name:         "комната",
		EmptyMessage: "пустая комната",
		EntryMessage: "ты в своей комнате",
		Containers:   []*models.Container{tableRoom, chair},
		IsEnabled:    true,
	}
	street := &models.Room{
		Name:          "улица",
		EmptyMessage:  "пустая улица",
		EntryMessage:  "на улице весна",
		EnableMessage: "дверь закрыта",
		Containers:    []*models.Container{door},
		IsEnabled:     false,
	}
	corridor := &models.Room{
		Name:           "коридор",
		EmptyMessage:   "пустой коридор",
		EntryMessage:   "ничего интересного",
		GlobalTransfer: "домой",
		Containers:     []*models.Container{door},
		IsEnabled:      true,
	}
	player := &models.Player{
		CurrentRoom: kitchen,
	}
	CurrentGame = models.NewGame(player, []*models.Room{room, corridor})
	CurrentGame.Actions["ключи дверь"] = func(player *models.Player) {
		if player.CurrentRoom.Name == "улица" {
			CurrentGame.Rooms["коридор"].IsEnabled = !CurrentGame.Rooms["коридор"].IsEnabled
		} else if player.CurrentRoom.Name == "коридор" {
			CurrentGame.Rooms["улица"].IsEnabled = !CurrentGame.Rooms["улица"].IsEnabled
		}
	}
	CurrentGame.Actions["надеть рюкзак"] = func(player *models.Player) {
		mission := CurrentGame.GetMission("собрать рюкзак")
		mission.IsCompleted = true
		player.CanTakeItems = true
	}
	CurrentGame.AddRoomTransfer(kitchen, corridor)
	CurrentGame.AddRoomTransfer(room, corridor)
	CurrentGame.AddGlobalTransfer(street, corridor)
	CurrentGame.SetMissions([]string{"собрать рюкзак", "идти в универ"})
}

func handleCommand(command string) string {
	splitCommand := strings.Split(command, " ")
	res := ""
	switch splitCommand[0] {
	case "осмотреться":
		res = CurrentGame.LookAround()
	case "идти":
		roomName := splitCommand[1]
		res = CurrentGame.Go(roomName)
	case "взять":
		itemName := splitCommand[1]
		res = CurrentGame.TakeItem(itemName)
	case "надеть":
		itemName := splitCommand[1]
		return CurrentGame.PutItem(itemName)
	case "применить":
		itemSubject := splitCommand[1]
		itemObject := splitCommand[2]
		res = CurrentGame.Apply(itemSubject, itemObject)
	default:
		res = "неизвестная команда"
	}
	return res
}
