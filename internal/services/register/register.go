package register

import (
	"errors"
	"fmt"
	"projectgrom/web/internal/repository/register"
)

type RegisterService struct {
	storage *register.StoragRegister
}

var (
	CreateError = errors.New("create register")
)

func NewRegisterService(data string) (*RegisterService, error) {
	reg, err := register.NewStorageRegister(data)
	if err != nil {
		return nil, fmt.Errorf("error: %s %w", err, CreateError)
	}
	return &RegisterService{storage: reg}, nil
}

func (r *RegisterService) Add(firstname, lastname, login, password string) error {
	return r.storage.Register(firstname, lastname, login, password)
}

func (r *RegisterService) UpdatePassword(login, password string) error {
	return r.storage.Update(login, password)
}

func (r *RegisterService) DeleteUser(login string) error {
	return r.storage.Delete(login)
}

func (r *RegisterService) GetUser(login, password string) error {
	return r.GetUser(login, password)
}
