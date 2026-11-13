package response

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/domain/dto"
)

func JSONSuccess(c fiber.Ctx, statusCode int, data any) error {
	return c.Status(statusCode).JSON(dto.Success(data))
}

func JSONError(c fiber.Ctx, statusCode int, code, message string) error {
	reqID, _ := c.Locals("requestid").(string)
	return c.Status(statusCode).JSON(dto.Fail(code, message, reqID))
}
