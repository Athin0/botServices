package repository

import (
	"botServices/pkg/model"
)

type IPasswordRepo interface {
	Set(owner int64, login string, password string, service string) (int, error)
	Get(owner int64, service string) (*model.ServiceInfo, error)
	Del(owner int64, service string) error
	GetAll(owner int64) ([]model.ServiceInfo, error)
}
