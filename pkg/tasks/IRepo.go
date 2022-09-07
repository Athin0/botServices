package tasks

import tgbotapi "github.com/skinass/telegram-bot-api/v5"

type ITaskRepo interface {
	ShowTasks(user *tgbotapi.Chat, bot *tgbotapi.BotAPI) error
	New(user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error
	Assign(user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error
	Unassign(user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error
	Resolve(user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error
	MyTasks(user *tgbotapi.Chat, bot *tgbotapi.BotAPI) error
	OwnTasks(user *tgbotapi.Chat, bot *tgbotapi.BotAPI) error
}