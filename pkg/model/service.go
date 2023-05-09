package model

import "fmt"

type ServiceInfo struct {
	ID        int
	Service   string
	Login     string
	Password  string
	ChatOwner int64
}

func (s *ServiceInfo) String() string {
	return fmt.Sprintf("Service:%s\nLogin:%s\nPassword:%s\n", s.Service, s.Login, s.Password)
}
