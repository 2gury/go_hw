package main

import (
	"context"
	"fmt"
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/tgbot/delivery"
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/tgbot/repository"
	"lectures-2022-1/04_net2/99_hw/taskbot/internal/tgbot/usecases"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "5138127632:AAHQLISNNgOdf5wyDUIIIKJB9CjpDW0NRy4"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://golang-hw4.herokuapp.com"
)

func startTaskBot(ctx context.Context) error {
	_, close := context.WithCancel(ctx)
	defer close()
	
	rep := repository.NewRepository()
	use := usecases.NewTgUsecase(rep)
	hnd := delivery.NewTgHandler(use, BotToken)

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("bot is working"))
		if err != nil {
			log.Println(err)
		}
	})

	bot, err := tgbotapi.NewBotAPI(BotToken)

	if err != nil {
		return err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		return fmt.Errorf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		return fmt.Errorf("SetWebhook failed: %s", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()

	updates := bot.ListenForWebhook("/")
	for update := range updates {
		if update.Message.IsCommand() {
			hnd.HandleMessageFromTg(update)
		}
	}

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		log.Println(err)
	}
}
