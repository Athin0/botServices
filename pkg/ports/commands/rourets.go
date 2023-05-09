package commands

import (
	"botServices/pkg/ports/handlers"
	"botServices/pkg/repository"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

const (
	Help   string = `help`
	Set           = `set`
	Get           = `get`
	Del           = `del`
	GetAll        = `getAll`
)

func Make(c repository.IPasswordRepo, command string, user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) (*tgbotapi.Message, error) {
	switch command {
	case Help:
		return nil, handlers.HandleHelp(user, bot)
	case Get:
		return handlers.HandleGet(c, user, bot, text)
	case Set:
		return handlers.HandleSet(c, user, bot, text)
	case Del:
		return handlers.HandleDel(c, user, bot, text)
	case GetAll:
		return handlers.HandleGetAll(c, user, bot, text)
	default:
		_, err := bot.Send(tgbotapi.NewMessage(user.ID, "Wrong command"))
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
