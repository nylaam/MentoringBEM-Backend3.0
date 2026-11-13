package order

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, handler *OrderHandler) {
	orders := router.Group("/orders")
	orders.Post("/", handler.Create)
	orders.Get("/", handler.GetAll)
	orders.Get("/:id", handler.GetByID)
	orders.Put("/:id", handler.Update)
	orders.Delete("/:id", handler.Delete)
}
