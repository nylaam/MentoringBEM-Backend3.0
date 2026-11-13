package order

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/domain/contracts"
	"github.com/mentoring/toko-backend/pkg/response"
)

type OrderHandler struct {
	svc contracts.OrderService
}

func NewOrderHandler(svc contracts.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) Create(c fiber.Ctx) error {
	// =========================================================================
	// TODO [Hands-on Step 3]: HTTP Delivery Layer (Thin Controller / Handler)
	// =========================================================================

	return fiber.NewError(fiber.StatusNotImplemented, "TODO [Hands-on Step 3]: implementasikan handler Create Order")
}

func (h *OrderHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.svc.GetOrderByID(c.Context(), id)
	if err != nil {
		return err
	}
	return response.JSONSuccess(c, fiber.StatusOK, res)
}

func (h *OrderHandler) GetAll(c fiber.Ctx) error {
	res, err := h.svc.GetAllOrders(c.Context())
	if err != nil {
		return err
	}
	return response.JSONSuccess(c, fiber.StatusOK, res)
}
