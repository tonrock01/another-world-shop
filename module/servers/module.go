package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tonrock01/another-world-shop/module/middlewares/middlewaresHandlers"
	"github.com/tonrock01/another-world-shop/module/middlewares/middlewaresRepositories"
	"github.com/tonrock01/another-world-shop/module/middlewares/middlewaresUsecases"
	"github.com/tonrock01/another-world-shop/module/monitor/monitorHandlers"
	"github.com/tonrock01/another-world-shop/module/users/usersHandlers"
	"github.com/tonrock01/another-world-shop/module/users/usersRepositories"
	"github.com/tonrock01/another-world-shop/module/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
}
