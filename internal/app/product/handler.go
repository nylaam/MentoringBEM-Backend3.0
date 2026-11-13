package product

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/domain/contracts"
	"github.com/mentoring/toko-backend/pkg/response"
)

type ProductHandler struct {
	svc contracts.ProductService
}

func NewProductHandler(svc contracts.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Create(c fiber.Ctx) error {
	// =========================================================================
	// TODO [Hands-on Step 3.1]: HTTP Delivery Layer - Create Product (POST)
	// =========================================================================
	// Sesuai prinsip Thin Handler (Slide 12 di index.md):
	// 1. Parsing & binding request body JSON ke dto.CreateProductRequest:
	// 2. Validasi format input via validator.ValidateStruct:
	// 3. Panggil domain service (Use Case)
	// 4. Return response terstandarisasi HTTP 201 Created:

	return fiber.NewError(fiber.StatusNotImplemented, "TODO [Hands-on Step 3.1]: implementasikan handler Create Product")
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
	// =========================================================================
	// TODO [Hands-on Step 3.2]: HTTP Delivery Layer - Update Product (PUT)
	// =========================================================================
	// 1. Ambil parameter ID dari URL route:
	// 2. Parsing request body JSON ke dto.UpdateProductRequest:
	// 3. Validasi format input via validator.ValidateStruct:
	// 4. Panggil domain service:
	// 5. Return response terstandarisasi HTTP 200 OK:

	return fiber.NewError(fiber.StatusNotImplemented, "TODO [Hands-on Step 3.2]: implementasikan handler Update Product")
}

func (h *ProductHandler) Delete(c fiber.Ctx) error {
	// =========================================================================
	// TODO [Hands-on Step 3.3]: HTTP Delivery Layer - Delete Product (DELETE)
	// =========================================================================
	// 1. Ambil parameter ID dari URL route:
	// 2. Panggil domain service untuk menghapus produk:
	// 3. Return response terstandarisasi HTTP 200 OK:

	return fiber.NewError(fiber.StatusNotImplemented, "TODO [Hands-on Step 3.3]: implementasikan handler Delete Product")
}
