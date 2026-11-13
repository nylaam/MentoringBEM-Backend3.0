package order

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/domain/contracts"
	"github.com/mentoring/toko-backend/domain/dto"
	"github.com/mentoring/toko-backend/pkg/response"
	"github.com/mentoring/toko-backend/pkg/validator"
)

type OrderHandler struct {
	svc contracts.OrderService
}

func NewOrderHandler(svc contracts.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) Create(c fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format JSON tidak valid")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	res, err := h.svc.CreateOrder(c.Context(), req)
	if err != nil {
		return err // Ditangani oleh centralized CustomErrorHandler di internal/middlewares
	}

	return response.JSONSuccess(c, fiber.StatusCreated, res)
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

func (h *OrderHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format JSON tidak valid")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	res, err := h.svc.UpdateOrder(c.Context(), id, req)
	if err != nil {
		return err
	}

	return response.JSONSuccess(c, fiber.StatusOK, res)
}

func (h *OrderHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.svc.DeleteOrder(c.Context(), id); err != nil {
		return err
	}

	return response.JSONSuccess(c, fiber.StatusOK, fiber.Map{
		"message": "Order berhasil dihapus",
	})
}
