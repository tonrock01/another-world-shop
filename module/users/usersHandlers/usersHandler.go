package usersHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/module/entities"
	"github.com/tonrock01/another-world-shop/module/users"
	"github.com/tonrock01/another-world-shop/module/users/usersUsecases"
)

type userHandlersErrCode string

const (
	signUpCustomerErr     userHandlersErrCode = "users-001"
	signInErr             userHandlersErrCode = "users-002"
	refreshPassportErr    userHandlersErrCode = "users-003"
	signOutErr            userHandlersErrCode = "users-004"
	signUpAdminErr        userHandlersErrCode = "users-005"
	generateAdminTokenErr userHandlersErrCode = "users-006"
	getUserProfileErr     userHandlersErrCode = "users-007"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg           config.IConfig
	usersUsecases usersUsecases.IUserUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecases usersUsecases.IUserUsecase) IUsersHandler {
	return &usersHandler{
		cfg:           cfg,
		usersUsecases: usersUsecases,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser

	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	// Insert
	result, err := h.usersUsecases.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
