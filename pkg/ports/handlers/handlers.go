package handlers

import (
	"botServices/pkg/repository"
	"errors"
	"fmt"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	"log"
	"strings"
)

func HandleHelp(user *tgbotapi.Chat, bot *tgbotapi.BotAPI) error {
	str := `
	/set LOGIN PASSWORD SERVICE     - добавляет данные о сервисе  (пробелы только между параметрами)
	/get	SERVICE	 				- получение данных для сервиса
	/del	SERVICE					- удаляет данные о сервисе
	/getAll      					- отображает все пароли
	/help        					- меню команд
	`
	_, err := bot.Send(tgbotapi.NewMessage(user.ID, str))
	return err
}

func HandleGet(c repository.IPasswordRepo, user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) (*tgbotapi.Message, error) {
	d := strings.Split(text, " ")
	if len(d) != 2 {
		return nil, errors.New("wrong info")
	}
	resp, err := c.Get(user.ID, d[1])
	msg, err := bot.Send(tgbotapi.NewMessage(user.ID, resp.String()))
	return &msg, err
}
func HandleSet(c repository.IPasswordRepo, user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) (*tgbotapi.Message, error) {
	d := strings.Split(text, " ")
	if len(d) != 4 {
		return nil, errors.New("don't enough info")
	}

	_, err := c.Set(user.ID, d[1], d[2], d[3])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	msg, err := bot.Send(tgbotapi.NewMessage(user.ID, "Successes set data!"))
	return &msg, err
}

func HandleDel(c repository.IPasswordRepo, user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) (*tgbotapi.Message, error) {
	d := strings.Split(text, " ")
	if len(d) != 2 {
		return nil, errors.New("don't enough info")
	}
	err := c.Del(user.ID, d[1])
	if err != nil {
		return nil, err
	}
	msg, err := bot.Send(tgbotapi.NewMessage(user.ID, "Successes delete data!"))
	return &msg, err
}
func HandleGetAll(c repository.IPasswordRepo, user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) (*tgbotapi.Message, error) {

	resp, err := c.GetAll(user.ID)
	var s string
	for i, info := range resp {
		s += fmt.Sprint(i+1) + "." + info.String()
	}
	if len(s) == 0 {
		_, err = bot.Send(tgbotapi.NewMessage(user.ID, "data empty"))
		return nil, err
	}
	msg, err := bot.Send(tgbotapi.NewMessage(user.ID, s))
	return &msg, err
}
