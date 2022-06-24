package tgbot

import "lectures-2022-1/04_net2/99_hw/taskbot/internal/models"

type TgRepository interface {
	SelectTasks() []*models.Task
	SelectUserByLogin(user models.User) *models.User
	InsertUser(user models.User) *models.User
	InsertTask(user *models.User, taskDescription string) *models.Task
	SelectTaskByID(user *models.User, taskID uint64) (int, *models.Task)
	RemoveTaskByID(taskID int)  
}