package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tonrock01/another-world-shop/modules/appinfo/appinfoHandlers"
	"github.com/tonrock01/another-world-shop/modules/appinfo/appinfoRepositories"
	"github.com/tonrock01/another-world-shop/modules/appinfo/appinfoUsecases"
	"github.com/tonrock01/another-world-shop/modules/files/filesHandlers"
	"github.com/tonrock01/another-world-shop/modules/files/filesUsecases"
	"github.com/tonrock01/another-world-shop/modules/middlewares/middlewaresHandlers"
	"github.com/tonrock01/another-world-shop/modules/middlewares/middlewaresRepositories"
	"github.com/tonrock01/another-world-shop/modules/middlewares/middlewaresUsecases"
	"github.com/tonrock01/another-world-shop/modules/monitor/monitorHandlers"
	"github.com/tonrock01/another-world-shop/modules/users/usersHandlers"
	"github.com/tonrock01/another-world-shop/modules/users/usersRepositories"
	"github.com/tonrock01/another-world-shop/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
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

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SingIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignUpAdmin)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoRepository(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) FilesModule() {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	router := m.r.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)
}
