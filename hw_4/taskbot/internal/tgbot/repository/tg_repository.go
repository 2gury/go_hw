package repository

import (
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/models"
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/tgbot"
	"sync"
)

type TgRepository struct {
	tasks      []*models.Task
	users      map[string]*models.User
	lastTaskID uint64

	mx *sync.Mutex
}

func NewRepository() tgbot.TgRepository {
	return &TgRepository{
		tasks:      []*models.Task{},
		users:      map[string]*models.User{},
		lastTaskID: 0,

		mx: &sync.Mutex{},
	}
}

func (r *TgRepository) GetLastTaskID() uint64 {
	r.lastTaskID++
	return r.lastTaskID
}

func (r *TgRepository) SelectTasks() []*models.Task {
	r.mx.Lock()
	defer r.mx.Unlock()

	return r.tasks
}

func (r *TgRepository) SelectUserByLogin(usr models.User) *models.User {
	r.mx.Lock()
	defer r.mx.Unlock()

	user, ok := r.users[usr.Login]
	if !ok {
		return nil
	}
	return user
}

func (r *TgRepository) InsertUser(usr models.User) *models.User {
	r.mx.Lock()
	defer r.mx.Unlock()

	user := &models.User{
		Login:  usr.Login,
		ChatID: usr.ChatID,
	}
	r.users[user.Login] = user
	return user
}

func (r *TgRepository) InsertTask(user *models.User, taskDescription string) *models.Task {
	r.mx.Lock()
	defer r.mx.Unlock()

	task := &models.Task{
		ID:          r.GetLastTaskID(),
		Description: taskDescription,
		CreatedBy:   user,
		AssignTo:    nil,
	}

	r.tasks = append(r.tasks, task)
	return task
}

func (r *TgRepository) RemoveTaskByID(taskID int) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if len(r.tasks) == 1 {
		r.tasks = r.tasks[:0]
	} else {
		r.tasks = append(r.tasks[:taskID], r.tasks[taskID+1:]...)
	}
}

func (r *TgRepository) SelectTaskByID(user *models.User, taskID uint64) (int, *models.Task) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for id, task := range r.tasks {
		if task.ID == taskID {
			return id, task
		}
	}
	return -1, nil
}
