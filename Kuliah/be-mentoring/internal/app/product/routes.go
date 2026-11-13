package product

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, handler *ProductHandler) {
	products := router.Group("/products")
	products.Post("/", handler.Create)
	products.Get("/", handler.GetAll)
	products.Get("/:id", handler.GetByID)
}
