package tgbot

import "lectures-2022-1/04_net2/99_hw/taskbot/internal/models"

type TgUsecase interface {
	GetTasks(user models.User) []*models.Output
	CreateTaskByUser(user models.User, messageText string) []*models.Output
	GetTasksForUser(user models.User) []*models.Output
	GetCreatedByUserTasks(user models.User) []*models.Output
	AssignTaskToUser(usr models.User, taskID uint64) []*models.Output
	UnassignTaskFromUser(usr models.User, taskID uint64) []*models.Output
	ResolveTask(usr models.User, taskID uint64) []*models.Output 
}
