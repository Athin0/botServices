package tasks

import (
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	"strconv"
)

const (
	NoTasks = "Нет задач"
)

type Task struct {
	ID         int
	Owner      string
	Task       string
	Worker     string
	chatWorker *tgbotapi.Chat
	chatOwner  *tgbotapi.Chat
}

type ArrayOfTasks struct {
	Tasks     []Task
	LastIndex int
}

func NewArrayOfTasks() *ArrayOfTasks {
	return &ArrayOfTasks{
		Tasks:     make([]Task, 0),
		LastIndex: 0,
	}

}

func (c *ArrayOfTasks) add(owner *tgbotapi.Chat, text string) int {
	c.LastIndex++

	task := Task{
		ID:        c.LastIndex,
		Owner:     owner.UserName,
		Task:      text,
		chatOwner: owner,
	}
	c.Tasks = append(c.Tasks, task)
	return task.ID
}

func (c *ArrayOfTasks) remove(id int) error {
	var index = -1

	for i, task := range c.Tasks {
		if id == task.ID {
			index = i
			break
		}
	}

	c.Tasks = append(c.Tasks[:index], c.Tasks[index+1:]...)
	return nil
}
func (c *ArrayOfTasks) find(id int) *Task {
	for i, task := range c.Tasks {
		if id == task.ID {
			return &c.Tasks[i]
		}
	}
	return nil
}

func (c *ArrayOfTasks) GetUserTasks(user string) ITaskRepo {
	arr := NewArrayOfTasks()
	for _, task := range c.Tasks {
		if task.Worker == user {
			arr.Tasks = append(arr.Tasks, task)
		}
	}
	return arr
}

func (c *ArrayOfTasks) GetOwnedTasks(user string) ITaskRepo {
	arr := NewArrayOfTasks()
	for _, task := range c.Tasks {
		if task.Owner == user {
			arr.Tasks = append(arr.Tasks, task)
		}
	}
	return arr
}

func (c *ArrayOfTasks) ShowTasks(user *tgbotapi.Chat, bot *tgbotapi.BotAPI) error {
	var str = ""
	for i, task := range c.Tasks {
		id := strconv.Itoa(task.ID)
		str += id + ". " + task.Task + " by @" + task.Owner + "\n"
		if task.Worker != "" {
			if task.Worker == user.UserName {
				str += "assignee: я" + "\n" +
					"/unassign_" + id + " /resolve_" + id
			} else {
				str += "assignee: " + "@" + task.Worker
			}
		} else {
			str += "/assign_" + id
		}
		if i != len(c.Tasks)-1 {
			str += "\n"
			str += "\n"
		}
	}
	if str == "" {
		str = NoTasks
	}
	_, err := bot.Send(tgbotapi.NewMessage(user.ID, str))
	return err
}

func (c *ArrayOfTasks) New(user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error {
	if text == "" {
		_, err0 := bot.Send(tgbotapi.NewMessage(user.ID, "задача не может быть пустой"))
		if err0 != nil {
			return err0
		}
		return nil
	}
	id := c.add(user, text)

	str := "Задача \"" + c.find(id).Task + "\" создана, id=" + strconv.Itoa(id)
	_, err0 := bot.Send(tgbotapi.NewMessage(user.ID, str))
	if err0 != nil {
		return err0
	}
	return nil
}

func (c *ArrayOfTasks) Assign(user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error {
	id, err := strconv.Atoi(text)
	if err != nil {
		return err
	}
	task := c.find(id)
	if task == nil {
		_, err := bot.Send(tgbotapi.NewMessage(user.ID, "Задачи с заданным ID не существует"))
		if err != nil {
			return err
		}
		return nil
	}
	prevWorker := task.chatWorker
	owner := task.chatOwner.ID

	if prevWorker != nil && prevWorker.ID == user.ID {
		_, err0 := bot.Send(tgbotapi.NewMessage(prevWorker.ID, "Задача уже на вас"))
		if err0 != nil {
			return err0
		}
		return nil
	}

	task.Worker = user.UserName
	task.chatWorker = user
	str := "Задача \"" + task.Task + "\" назначена на вас"
	str1 := "Задача \"" + task.Task + "\" назначена на @" + user.UserName

	_, err1 := bot.Send(tgbotapi.NewMessage(user.ID, str))
	if err1 != nil {
		return err1
	}
	if prevWorker != nil && prevWorker.ID != user.ID {
		_, err0 := bot.Send(tgbotapi.NewMessage(prevWorker.ID, str1))
		if err0 != nil {
			return err0
		}
		return nil
	}
	if owner != user.ID {
		_, err0 := bot.Send(tgbotapi.NewMessage(owner, str1))
		if err0 != nil {
			return err0
		}
		return nil
	}
	return nil
}
func (c *ArrayOfTasks) Unassign(user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error {
	var str string
	id, err := strconv.Atoi(text)
	if err != nil {
		return err
	}
	task := c.find(id)
	if task == nil {
		_, err := bot.Send(tgbotapi.NewMessage(user.ID, "Задачи с заданным ID не существует"))
		if err != nil {
			return err
		}
		return nil
	}
	owner := task.chatOwner.ID
	if task.Worker == user.UserName {
		str = `Принято`
		task.Worker = ""
		task.chatWorker = nil
		str1 := `Задача "` + task.Task + `" осталась без исполнителя`
		if owner != user.ID {
			_, err0 := bot.Send(tgbotapi.NewMessage(owner, str1))
			if err0 != nil {
				return err0
			}
		}
	} else {
		str = `Задача не на вас`
	}

	_, err1 := bot.Send(tgbotapi.NewMessage(user.ID, str))
	if err1 != nil {
		return err1
	}
	return nil
}

func (c *ArrayOfTasks) Resolve(user *tgbotapi.Chat, bot *tgbotapi.BotAPI, text string) error {
	var str string
	var str1 string
	id, err := strconv.Atoi(text)
	if err != nil {
		return err
	}
	task := c.find(id)
	if task == nil {
		_, err := bot.Send(tgbotapi.NewMessage(user.ID, "Задачи с заданным ID не существует"))
		if err != nil {
			return err
		}
		return nil
	}
	owner := task.chatOwner.ID
	if task.Worker == user.UserName || task.Owner == user.UserName {
		str = `Задача "` + task.Task + `" выполнена`
		str1 = `Задача "` + task.Task + `" выполнена @` + user.UserName
		err := c.remove(id)
		if err != nil {
			return err
		}
		if owner != user.ID {
			_, err0 := bot.Send(tgbotapi.NewMessage(owner, str1))
			if err0 != nil {
				return err0
			}
		}
	} else {
		str = `Задача не на вас, чтобы выполнить используйте: /assign_` + strconv.Itoa(task.ID)
	}
	_, err1 := bot.Send(tgbotapi.NewMessage(user.ID, str))
	if err1 != nil {
		return err1
	}
	return nil

}

func (c *ArrayOfTasks) MyTasks(user *tgbotapi.Chat, bot *tgbotapi.BotAPI) error {
	var str string
	for i, task := range c.Tasks {
		id := strconv.Itoa(task.ID)
		if task.Worker == user.UserName {
			if i != len(c.Tasks) && str != "" {
				str += "\n"
				str += "\n"
			}
			str += id + ". " + task.Task + " by @" + task.Owner + "\n"
			str += "/unassign_" + id + " /resolve_" + id
		}
	}
	if str == "" {
		str = NoTasks
	}
	_, err1 := bot.Send(tgbotapi.NewMessage(user.ID, str))
	if err1 != nil {
		return err1
	}
	return nil
}

func (c *ArrayOfTasks) OwnTasks(user *tgbotapi.Chat, bot *tgbotapi.BotAPI) error {
	var str string
	for i, task := range c.Tasks {
		id := strconv.Itoa(task.ID)
		if task.Owner == user.UserName {
			if i != len(c.Tasks) && str != "" {
				str += "\n" + "\n"
			}
			str += id + ". " + task.Task + " by @" + task.Owner + "\n"
			if task.Worker != "" {
				if task.Worker == user.UserName {
					str += "/unassign_" + id + " /resolve_" + id
				}
			} else {
				str += "/assign_" + id
			}
		}
	}
	if str == "" {
		str = NoTasks
	}
	_, err1 := bot.Send(tgbotapi.NewMessage(user.ID, str))
	if err1 != nil {
		return err1
	}
	return nil
}
