package v1

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	models "github.com/bulatok/ozon-task/internal/ozon-task/models/api/v1"
)

func (a *ApiServer) initLinksRouters(r fiber.Router) {
	{
		r.Post("/new", a.newLink())
		r.Get("/:linkHash", a.getLink())
	}
}

// newLink
// POST /links/new
func (a *ApiServer) newLink() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		a.l.Debug("new link", zap.String("path", ctx.Path()))

		reqBody := &models.ApiNewLinkRequest{}
		if err := json.Unmarshal(ctx.Body(), reqBody); err != nil {
			a.l.Info("invalid request body", zap.String("error", err.Error()))
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		shortLink, err := a.uc.Links.New(reqBody.OriginalLink, a.conf.Service.PublicUrl)
		if err != nil {
			return toFiberError(fiber.StatusInternalServerError, err)
		}

		return ctx.JSON(models.ApiNewLinkResponse{ShortLink: shortLink})
	}
}

// getLink
// GET /links/:linkHash
func (a *ApiServer) getLink() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		a.l.Debug("get link", zap.String("path", ctx.Path()))

		linkHash := ctx.Params("linkHash")
		if linkHash == "" {
			a.l.Info("invalid link hash")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		original, err := a.uc.Links.Get(a.conf.Service.PublicUrl + "/" + linkHash)
		if err != nil {
			return toFiberError(fiber.StatusInternalServerError, err)
		}

		return ctx.Status(fiber.StatusOK).JSON(models.ApiGetOriginalLinkResponse{OriginalLink: original})
	}
}
