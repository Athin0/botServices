package repository

import (
	"botServices/pkg/model"
)

const (
	NoService = "Нет задач"
)

type ArrayOfService struct {
	Service   map[string]model.ServiceInfo
	LastIndex int
}

func NewArrayOfService() *ArrayOfService {
	return &ArrayOfService{
		Service:   make(map[string]model.ServiceInfo),
		LastIndex: 0,
	}

}

func (c *ArrayOfService) Set(owner int64, login string, password string, service string) (int, error) {
	c.LastIndex++
	task := model.ServiceInfo{
		ID:        c.LastIndex,
		Login:     login,
		Password:  password,
		Service:   service,
		ChatOwner: owner,
	}
	c.Service[service] = task
	return task.ID, nil
}

func (c *ArrayOfService) Del(owner int64, service string) error {
	delete(c.Service, service)
	return nil
}

func (c *ArrayOfService) Get(owner int64, service string) (*model.ServiceInfo, error) {
	m := (c.Service[service])
	return &m, nil
}

func (c *ArrayOfService) GetAll(owner int64) ([]model.ServiceInfo, error) {
	arr := []model.ServiceInfo{}
	for _, task := range c.Service {
		if task.ChatOwner == owner {
			arr = append(arr, task)
		}
	}
	return arr, nil
}
