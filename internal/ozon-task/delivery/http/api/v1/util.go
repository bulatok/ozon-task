package v1

import (
	"github.com/bulatok/ozon-task/internal/ozon-task/models"
	"github.com/gofiber/fiber/v2"
)

type ErrorData struct {
	Value string `json:"value"`
}

type errorResponse struct {
	Error bool      `json:"error"`
	Data  ErrorData `json:"data"`
}

func toFiberError(statusCode int, err error) *fiber.Error {
	commonErr, ok := err.(models.CommonError)
	if ok {
		return fiber.NewError(commonErr.StatusCode, err.Error())
	}
	return fiber.NewError(statusCode, err.Error())
}
