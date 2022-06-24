package usecases

import (
	"fmt"
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/models"
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/tgbot"
)

const (
	NoTasks = "Нет задач"
)

type TgUsecase struct {
	tgRep tgbot.TgRepository
}

func NewTgUsecase(rep tgbot.TgRepository) tgbot.TgUsecase {
	return &TgUsecase{
		tgRep: rep,
	}
}

func (u *TgUsecase) GetTasks(usr models.User) []*models.Output {
	tasks := u.tgRep.SelectTasks()
	user := u.tgRep.SelectUserByLogin(usr)
	if user == nil {
		user = u.tgRep.InsertUser(usr)
	}
	resStr := ""

	for id, task := range tasks {
		resStr += fmt.Sprintf("%d. %s by %s", task.ID, task.Description, task.CreatedBy.Login)
		switch {
		case task.AssignTo == user:
			resStr += fmt.Sprintf("\nassignee: я\n/unassign_%d /resolve_%d", task.ID, task.ID)
		case task.AssignTo != nil:
			resStr += fmt.Sprintf("\nassignee: %s", task.AssignTo.Login)
		default:
			resStr += fmt.Sprintf("\n/assign_%d", task.ID)
		}
		if id != len(tasks)-1 {
			resStr += "\n\n"
		}
	}

	if resStr == "" {
		resStr = NoTasks
	}
	return []*models.Output{
		{
			Message: resStr,
			ChatID:  user.ChatID,
		},
	}
}

func (u *TgUsecase) GetTasksForUser(usr models.User) []*models.Output {
	tasks := u.tgRep.SelectTasks()
	user := u.tgRep.SelectUserByLogin(usr)
	if user == nil {
		user = u.tgRep.InsertUser(usr)
	}
	resStr := ""

	userTasks := []*models.Task{}
	for _, task := range tasks {
		if task.AssignTo == user {
			userTasks = append(userTasks, task)
		}
	}
	for id, task := range userTasks {
		if task.AssignTo == user {
			resStr += fmt.Sprintf("%d. %s by %s", task.ID, task.Description, task.CreatedBy.Login)
			resStr += fmt.Sprintf("\n/unassign_%d /resolve_%d", task.ID, task.ID)
			if id != len(userTasks)-1 {
				resStr += "\n\n"
			}
		}
	}

	if resStr == "" {
		resStr = NoTasks
	}
	return []*models.Output{
		{
			Message: resStr,
			ChatID:  user.ChatID,
		},
	}
}

func (u *TgUsecase) GetCreatedByUserTasks(usr models.User) []*models.Output {
	tasks := u.tgRep.SelectTasks()
	user := u.tgRep.SelectUserByLogin(usr)
	if user == nil {
		user = u.tgRep.InsertUser(usr)
	}
	resStr := ""
	userTasks := []*models.Task{}
	for _, task := range tasks {
		if task.CreatedBy == user {
			userTasks = append(userTasks, task)
		}
	}
	for id, task := range userTasks {
		if task.CreatedBy == user {
			resStr += fmt.Sprintf("%d. %s by %s", task.ID, task.Description, task.CreatedBy.Login)
			switch {
			case task.AssignTo == user:
				resStr += fmt.Sprintf("\n/unassign_%d /resolve_%d", task.ID, task.ID)
			default:
				resStr += fmt.Sprintf("\n/assign_%d", task.ID)
			}
			if id != len(userTasks)-1 {
				resStr += "\n\n"
			}
		}
	}

	if resStr == "" {
		resStr = NoTasks
	}
	return []*models.Output{
		{
			Message: resStr,
			ChatID:  user.ChatID,
		},
	}
}

func (u *TgUsecase) CreateTaskByUser(usr models.User, taskDescription string) []*models.Output {
	user := u.tgRep.SelectUserByLogin(usr)
	if user == nil {
		user = u.tgRep.InsertUser(usr)
	}
	task := u.tgRep.InsertTask(user, taskDescription)

	return []*models.Output{
		{
			Message: fmt.Sprintf(`Задача "%s" создана, id=%d`, task.Description, task.ID),
			ChatID:  user.ChatID,
		},
	}
}

func (u *TgUsecase) AssignTaskToUser(usr models.User, taskID uint64) []*models.Output {
	user := u.tgRep.SelectUserByLogin(usr)
	var output []*models.Output
	if user == nil {
		user = u.tgRep.InsertUser(usr)
	}
	_, task := u.tgRep.SelectTaskByID(user, taskID)
	if task == nil {
		output = []*models.Output{
			{
				Message: "Такой задачи нет",
				ChatID:  user.ChatID,
			},
		}
		return output
	}
	lastAssigner := task.AssignTo
	creator := task.CreatedBy
	task.AssignTo = user

	switch {
	case user == creator:
		output = []*models.Output{
			{
				Message: fmt.Sprintf(`Задача "%s" назначена на вас`, task.Description),
				ChatID:  user.ChatID,
			},
		}
	case lastAssigner == nil:
		output = []*models.Output{
			{
				Message: fmt.Sprintf(`Задача "%s" назначена на вас`, task.Description),
				ChatID:  user.ChatID,
			},
			{
				Message: fmt.Sprintf(`Задача "%s" назначена на %s`, task.Description, user.Login),
				ChatID:  creator.ChatID,
			},
		}
	default:
		output = []*models.Output{
			{
				Message: fmt.Sprintf(`Задача "%s" назначена на вас`, task.Description),
				ChatID:  user.ChatID,
			},
			{
				Message: fmt.Sprintf(`Задача "%s" назначена на %s`, task.Description, user.Login),
				ChatID:  lastAssigner.ChatID,
			},
		}
	}

	return output
}

func (u *TgUsecase) UnassignTaskFromUser(usr models.User, taskID uint64) []*models.Output {
	user := u.tgRep.SelectUserByLogin(usr)
	var output []*models.Output
	if user == nil {
		user = u.tgRep.InsertUser(usr)
	}
	_, task := u.tgRep.SelectTaskByID(user, taskID)
	if task == nil {
		output = []*models.Output{
			{
				Message: "Такой задачи нет",
				ChatID:  user.ChatID,
			},
		}
		return output
	}
	switch {
	case user != task.AssignTo:
		output = []*models.Output{
			{
				Message: "Задача не на вас",
				ChatID:  user.ChatID,
			},
		}
	case user != task.CreatedBy:
		output = []*models.Output{
			{
				Message: `Принято`,
				ChatID:  user.ChatID,
			},
			{
				Message: fmt.Sprintf(`Задача "%s" осталась без исполнителя`, task.Description),
				ChatID:  task.CreatedBy.ChatID,
			},
		}
		task.AssignTo = nil
	case user == task.CreatedBy:
		output = []*models.Output{
			{
				Message: "Принято",
				ChatID:  user.ChatID,
			},
		}
		task.AssignTo = nil
	}

	return output
}

func (u *TgUsecase) ResolveTask(usr models.User, taskID uint64) []*models.Output {
	user := u.tgRep.SelectUserByLogin(usr)
	var output []*models.Output
	if user == nil {
		user = u.tgRep.InsertUser(usr)
	}
	taskIndex, task := u.tgRep.SelectTaskByID(user, taskID)
	if task == nil {
		output = []*models.Output{
			{
				Message: "Такой задачи нет",
				ChatID:  user.ChatID,
			},
		}
		return output
	}

	switch {
	case user != task.AssignTo:
		output = []*models.Output{
			{
				Message: "Задача не на вас",
				ChatID:  user.ChatID,
			},
		}
	case user != task.CreatedBy:
		output = []*models.Output{
			{
				Message: fmt.Sprintf(`Задача "%s" выполнена`, task.Description),
				ChatID:  user.ChatID,
			},
			{
				Message: fmt.Sprintf(`Задача "%s" выполнена %s`, task.Description, task.AssignTo.Login),
				ChatID:  task.CreatedBy.ChatID,
			},
		}
		u.tgRep.RemoveTaskByID(taskIndex)
	case user == task.CreatedBy:
		output = []*models.Output{
			{
				Message: fmt.Sprintf(`Задача "%s" выполнена`, task.Description),
				ChatID:  user.ChatID,
			},
		}
		u.tgRep.RemoveTaskByID(taskIndex)
	}

	return output
}
