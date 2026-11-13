package order

import (
	"context"
	"fmt"

	"github.com/mentoring/toko-backend/domain/contracts"
	"github.com/mentoring/toko-backend/domain/dto"
	"github.com/mentoring/toko-backend/domain/entity"
)

type orderService struct {
	orderRepo   contracts.OrderRepository
	productRepo contracts.ProductRepository
}

func NewOrderService(orderRepo contracts.OrderRepository, productRepo contracts.ProductRepository) contracts.OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// 1. Ambil data produk & transform DTO ke domain items
	var domainItems []entity.OrderItem
	for _, it := range req.Items {
		prod, err := s.productRepo.FindByID(ctx, it.ProductID)
		if err != nil {
			return nil, err
		}

		domainItems = append(domainItems, entity.OrderItem{
			ProductID:   prod.ID,
			ProductName: prod.Name,
			Quantity:    it.Quantity,
			UnitPrice:   prod.Price,
		})
	}

	// 2. Eksekusi aturan bisnis via domain constructor
	order, err := entity.NewOrder(req.CustomerID, req.CustomerName, domainItems)
	if err != nil {
		return nil, err
	}

	// 3. Simpan order ke database via repository interface
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("gagal menyimpan pesanan: %w", err)
	}

	// 4. Return response DTO
	return toOrderResponse(order), nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (*dto.OrderResponse, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toOrderResponse(o), nil
}

func (s *orderService) GetAllOrders(ctx context.Context) ([]*dto.OrderResponse, error) {
	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []*dto.OrderResponse
	for _, o := range orders {
		res = append(res, toOrderResponse(o))
	}
	return res, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, id string, req dto.UpdateOrderRequest) (*dto.OrderResponse, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := o.UpdateStatus(req.Status); err != nil {
		return nil, err
	}

	if req.CustomerName != "" {
		o.CustomerName = req.CustomerName
	}

	if err := s.orderRepo.Update(ctx, o); err != nil {
		return nil, fmt.Errorf("gagal mengupdate order: %w", err)
	}

	return toOrderResponse(o), nil
}

func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	_, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.orderRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("gagal menghapus order: %w", err)
	}

	return nil
}

func toOrderResponse(o *entity.Order) *dto.OrderResponse {
	var itemsRes []dto.OrderItemResponse
	for _, it := range o.Items {
		itemsRes = append(itemsRes, dto.OrderItemResponse{
			ProductID:   it.ProductID,
			ProductName: it.ProductName,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
			Subtotal:    it.Subtotal,
		})
	}

	return &dto.OrderResponse{
		ID:           o.ID,
		CustomerID:   o.CustomerID,
		CustomerName: o.CustomerName,
		Items:        itemsRes,
		Total:        o.Total,
		Status:       o.Status,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}
