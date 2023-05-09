package main

import (
	"botServices/pkg/ports/commands"
	"fmt"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"time"

	"botServices/db"
	"botServices/pkg/repository"
	"botServices/secret"
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
	list, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

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

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "80"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	delChan := make(chan tgbotapi.DeleteMessageConfig)
	defer close(delChan)
	go func(delchan chan tgbotapi.DeleteMessageConfig) {
		for r := range delchan {
			_, err := bot.Request(r)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}(delChan)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		go handleUpdate(bot, update, list, delChan)

	}

	return nil
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, list repository.IPasswordRepo, delChan chan tgbotapi.DeleteMessageConfig) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic: %s", p)
		}
	}()

	command := update.Message.Command()
	text := update.Message.Text
	if command == "start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome!")
		msg.ReplyMarkup = keyBoard
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
		command = "help"
	}
	user := update.Message.Chat

	message, err := commands.Make(list, command, user, bot, text)
	if err != nil {
		log.Println(err.Error())
	}
	if message != nil {
		<-time.After(10 * time.Second)
		delChan <- tgbotapi.NewDeleteMessage(user.ID, message.MessageID)
	}
}

var keyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/help"),
		tgbotapi.NewKeyboardButton("/getAll"),
	),
)
