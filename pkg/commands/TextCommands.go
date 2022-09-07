package commands

import (
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	"strconv"
	"strings"
	"tgBotTasks/pkg/tasks"
)

const (
	Help     string = `/help`
	Tasks           = `/tasks`
	New             = `/new`
	Assign          = `/assign_`
	Unassign        = `/unassign_`
	Resolve         = `/resolve_`
	My              = `/my`
	Owner           = `/owner`
)

func HelpList(user *tgbotapi.Chat, bot *tgbotapi.BotAPI) error {
	str := `
	/tasks        -отображает существующие задачи
	/new XXX      - создаёт новую задачу
	/assign_$ID   - сделать себя исполнителем задачи №ID
	/unassign_$ID - снимает задачу с вас
	/resolve_$ID  - выполняет задачу, удаляет её из списка
	/my           - показывает задачи, которые назначены на меня
	/owner        - показывает задачи которые были созданы мной
	/help         - меню команд
	`
	_, err := bot.Send(tgbotapi.NewMessage(user.ID, str))
	return err
}

func SelectCommand(text string) (string, bool, string) {

	if len(text) > 5 && text[0:6] == `/start` {
		return Help, true, ""
	}
	if len(text) > 4 && text[0:5] == "/new " {
		return New, true, strings.Split(text, "/new ")[1]
	}
	if len(text) > 7 && text[0:8] == "/assign_" {
		id, err := strconv.Atoi(strings.Split(text, "_")[1])
		if err != nil {
			return Help, false, ""
		}
		return Assign, true, strconv.Itoa(id)
	}
	if len(text) > 9 && text[0:10] == "/unassign_" {
		id, err := strconv.Atoi(strings.Split(text, "_")[1])
		if err != nil {
			return Help, false, ""
		}
		return Unassign, true, strconv.Itoa(id)
	}
	if len(text) > 8 && text[0:9] == "/resolve_" {
		id, err := strconv.Atoi(strings.Split(text, "_")[1])
		if err != nil {
			return Help, false, ""
		}
		return Resolve, true, strconv.Itoa(id)
	}

	return text, true, ""
}

func Make(c tasks.ITaskRepo, command string, user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error {
	switch command {
	case Help:
		return HelpList(user, bot)
	case Tasks:
		return c.ShowTasks(user, bot)
	case New:
		return c.New(user, bot, text)
	case Assign:
		return c.Assign(user, bot, text)
	case Unassign:
		return c.Unassign(user, bot, text)
	case Resolve:
		return c.Resolve(user, bot, text)
	case My:
		return c.MyTasks(user, bot)
	case Owner:
		return c.OwnTasks(user, bot)

	}
	return nil
}
