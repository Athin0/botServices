package main

import (
	"fmt"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	"log"
	"math/rand"
	"net/http"
	"os"
	"tgBotTasks/pkg/commands"
	"tgBotTasks/pkg/tasks"
	"tgBotTasks/secret"
	"time"
)

var (
	BotToken   = secret.BotToken
	WebhookURL = secret.WebhookURL
)

func main() {
	err := startTaskBot()
	if err != nil {
		panic(err)
	}
}
func startTaskBot() error {
	rand.Seed(time.Now().UnixNano())
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
		return err
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
		return err
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
		return err
	}

	updates := bot.ListenForWebhook("/")

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("all is working"))
		if err != nil {
			log.Println("find state err: ", err.Error())
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	list := tasks.NewArrayOfTasks()

	for update := range updates {
		handleUpdate(bot, update, list)
	}
	return nil
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, list tasks.ITaskRepo) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic: %s", p)
		}
	}()

	if update.Message == nil {
		return
	}
	log.Printf("text: %#v\n", update.Message.Text)
	log.Printf("upd: %#v\n", update)

	command, ok, text := commands.SelectCommand(update.Message.Text)
	if !ok {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			`Неизвестная команда`,
		)

		msg.ReplyMarkup = &tgbotapi.ReplyKeyboardMarkup{
			Keyboard: [][]tgbotapi.KeyboardButton{
				{
					{
						Text: "/help",
					},
				},
			},
		}
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err.Error())
		}
		return
	}
	user := update.Message.Chat

	err := commands.Make(list, command, user, bot, text)

	if err != nil {
		log.Println(err.Error())
	}
}
