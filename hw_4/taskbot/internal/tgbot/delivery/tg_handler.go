package delivery

import (
	"fmt"
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/models"
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/tgbot"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

type TgHandler struct {
	tgUse    tgbot.TgUsecase
	botToken string
}

func NewTgHandler(use tgbot.TgUsecase, tgBotoken string) *TgHandler {
	return &TgHandler{
		tgUse:    use,
		botToken: tgBotoken,
	}
}

func (h *TgHandler) HandleMessageFromTg(update tgbotapi.Update) {
	type MessageBody struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}

	outputMessages := h.HandleCommandMessage(update)

	log.Printf("\n%d\n", len(outputMessages))
	for _, msg := range outputMessages {
		params := url.Values{
			"chat_id": {fmt.Sprintf("%d", msg.ChatID)},
			"text":    {msg.Message},
		}

		reqURL := fmt.Sprintf(tgbotapi.APIEndpoint, h.botToken, "sendMessage")

		log.Println(reqURL)

		_, err := http.Post(reqURL, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (h *TgHandler) HandleCommandMessage(update tgbotapi.Update) []*models.Output {
	output := []*models.Output{}

	user := models.User{
		Login:  "@" + update.Message.From.UserName,
		ChatID: update.Message.Chat.ID,
	}

	msgCommand := update.Message.Command()
	switch {
	case msgCommand == "tasks":
		output = h.GetTasks(user)
	case msgCommand == "new":
		output = h.CreateTask(user, update.Message.Text)
	case msgCommand == "my":
		output = h.GetUserTasks(user)
	case msgCommand == "owner":
		output = h.GetCreatedByUserTasks(user)
	case strings.HasPrefix(msgCommand, "unassign_"):
		output = h.UnassignTaskFromUser(user, update.Message.Text)
	case strings.HasPrefix(msgCommand, "assign_"):
		output = h.AssignTaskToUser(user, update.Message.Text)

	case strings.HasPrefix(msgCommand, "resolve_"):
		output = h.ResolveTask(user, update.Message.Text)
	default:
		output = append(output, &models.Output{
			Message: "unknown command",
			ChatID:  user.ChatID,
		})
	}

	return output
}

func (h *TgHandler) ResolveTask(user models.User, msgText string) []*models.Output {
	taskDescription := strings.ReplaceAll(msgText, "/resolve_", "")
	taskDescription = strings.TrimSpace(taskDescription)
	if taskDescription != "" {
		taskID, err := strconv.Atoi(taskDescription)
		if err != nil {
			return []*models.Output{
				{
					Message: "Укажите корректный номер задачи",
					ChatID:  user.ChatID,
				},
			}
		}
		output := h.tgUse.ResolveTask(user, uint64(taskID))
		return output
	}

	return []*models.Output{
		{
			Message: "Укажите корректный номер задачи",
			ChatID:  user.ChatID,
		},
	}
}

func (h *TgHandler) AssignTaskToUser(user models.User, msgText string) []*models.Output {
	taskDescription := strings.ReplaceAll(msgText, "/assign_", "")
	taskDescription = strings.TrimSpace(taskDescription)
	if taskDescription != "" {
		taskID, err := strconv.Atoi(taskDescription)
		if err != nil {
			return []*models.Output{
				{
					Message: "Укажите корректный номер задачи",
					ChatID:  user.ChatID,
				},
			}
		}
		output := h.tgUse.AssignTaskToUser(user, uint64(taskID))
		return output
	}

	return []*models.Output{
		{
			Message: "Укажите корректный номер задачи",
			ChatID:  user.ChatID,
		},
	}
}

func (h *TgHandler) UnassignTaskFromUser(user models.User, msgText string) []*models.Output {
	taskDescription := strings.ReplaceAll(msgText, "/unassign_", "")
	taskDescription = strings.TrimSpace(taskDescription)
	if taskDescription != "" {
		taskID, err := strconv.Atoi(taskDescription)
		if err != nil {
			return []*models.Output{
				{
					Message: "Укажите корректный номер задачи",
					ChatID:  user.ChatID,
				},
			}
		}
		output := h.tgUse.UnassignTaskFromUser(user, uint64(taskID))
		return output
	}

	return []*models.Output{
		{
			Message: "Укажите корректный номер задачи",
			ChatID:  user.ChatID,
		},
	}
}

func (h *TgHandler) GetTasks(user models.User) []*models.Output {
	return h.tgUse.GetTasks(user)
}

func (h *TgHandler) GetUserTasks(user models.User) []*models.Output {
	return h.tgUse.GetTasksForUser(user)
}

func (h *TgHandler) GetCreatedByUserTasks(user models.User) []*models.Output {
	return h.tgUse.GetCreatedByUserTasks(user)
}

func (h *TgHandler) CreateTask(user models.User, msgText string) []*models.Output {
	taskDescription := strings.ReplaceAll(msgText, "/new", "")
	taskDescription = strings.TrimSpace(taskDescription)
	if taskDescription != "" {
		return h.tgUse.CreateTaskByUser(user, taskDescription)
	}

	return []*models.Output{
		{
			Message: "Укажите название задачи",
			ChatID:  user.ChatID,
		},
	}
}
