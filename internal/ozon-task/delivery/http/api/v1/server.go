package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"

	"github.com/bulatok/ozon-task/internal/ozon-task/config"
	"github.com/bulatok/ozon-task/internal/ozon-task/models"
	"github.com/bulatok/ozon-task/internal/ozon-task/usecase"
)

var (
	ErrorHandler = func(ctx *fiber.Ctx, err error) error {
		c := fiber.StatusBadRequest
		if e, ok := err.(*fiber.Error); ok {
			c = e.Code
		}

		commonErr, ok := err.(models.CommonError)
		if ok {
			c = commonErr.StatusCode
		}

		errResp := errorResponse{
			Error: true,
			Data: ErrorData{
				Value: err.Error(),
			},
		}
		return ctx.Status(c).JSON(errResp)
	}
)

type ApiServer struct {
	uc   UseCases
	conf *config.Config
	l    *zap.Logger
}

type UseCases struct {
	Links *usecase.Links
}

func ProvideServer(conf *config.Config, l *zap.Logger, apiConfig UseCases) *ApiServer {
	childL := l.With(zap.String("logger", "http server"))
	return &ApiServer{
		conf: conf,
		l:    childL,
		uc:   apiConfig,
	}
}

func (a *ApiServer) Init(app *fiber.App) {
	app.Use(cors.New())

	// init routers
	rMain := app.Group("")

	rMain.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("pong"))
	})

	a.initLinksRouters(rMain)

}
