package middlewares

import (
	stderrs "errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/domain/dto"
	"github.com/mentoring/toko-backend/domain/errx"
)

// CustomErrorHandler menangani error terpusat di Fiber v3 (Slide 13 di index.md)
func CustomErrorHandler(c fiber.Ctx, err error) error {
	reqID, _ := c.Locals("requestid").(string)

	statusCode := http.StatusInternalServerError
	errorCode := "INTERNAL_SERVER_ERROR"
	errorMessage := "Terjadi kesalahan internal server"

	if domainErr, ok := errx.As(err); ok {
		errorCode = domainErr.Code
		errorMessage = domainErr.Message

		switch domainErr.Kind {
		case errx.KindNotFound:
			statusCode = http.StatusNotFound
		case errx.KindInvalid:
			statusCode = http.StatusUnprocessableEntity
		case errx.KindConflict:
			statusCode = http.StatusConflict
		case errx.KindUnauthorized:
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}
	} else {
		var fiberErr *fiber.Error
		if stderrs.As(err, &fiberErr) {
			statusCode = fiberErr.Code
			errorCode = http.StatusText(fiberErr.Code)
			errorMessage = fiberErr.Message
		}
	}

	return c.Status(statusCode).JSON(dto.Fail(errorCode, errorMessage, reqID))
}
