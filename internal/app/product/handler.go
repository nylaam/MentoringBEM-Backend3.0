package product

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/domain/contracts"
	"github.com/mentoring/toko-backend/domain/dto"
	"github.com/mentoring/toko-backend/pkg/response"
	"github.com/mentoring/toko-backend/pkg/validator"
)

type ProductHandler struct {
	svc contracts.ProductService
}

func NewProductHandler(svc contracts.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Create(c fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format JSON tidak valid")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	res, err := h.svc.CreateProduct(c.Context(), req)
	if err != nil {
		return err
	}

	return response.JSONSuccess(c, fiber.StatusCreated, res)
}

func (h *ProductHandler) GetAll(c fiber.Ctx) error {
	res, err := h.svc.GetAllProducts(c.Context())
	if err != nil {
		return err
	}
	return response.JSONSuccess(c, fiber.StatusOK, res)
}

func (h *ProductHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.svc.GetProductByID(c.Context(), id)
	if err != nil {
		return err
	}
	return response.JSONSuccess(c, fiber.StatusOK, res)
}

func (h *ProductHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format JSON tidak valid")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	res, err := h.svc.UpdateProduct(c.Context(), id, req)
	if err != nil {
		return err
	}

	return response.JSONSuccess(c, fiber.StatusOK, res)
}

func (h *ProductHandler) Patch(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.PatchProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format JSON tidak valid")
	}

	// Minimal satu field harus disediakan untuk operasi patch
	if req.Name == nil && req.Price == nil && req.Stock == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Minimal satu field (name, price, stock) harus diisi untuk patch")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	res, err := h.svc.PatchProduct(c.Context(), id, req)
	if err != nil {
		return err
	}

	return response.JSONSuccess(c, fiber.StatusOK, res)
}

func (h *ProductHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.svc.DeleteProduct(c.Context(), id); err != nil {
		return err
	}

	return response.JSONSuccess(c, fiber.StatusOK, fiber.Map{
		"message": "Produk berhasil dihapus",
	})
}
