package usersUsecases

import (
	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/module/users"
	"github.com/tonrock01/another-world-shop/module/users/usersRepositories"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUserUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	//Hashing a password
	if err := req.BcryotHashing(); err != nil {
		return nil, err
	}

	//Insert user
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}
