package repository

import (
	"context"
	"gorm.io/gorm"
	"order-service/infrastructure/log"
	"order-service/internal/entity"
)

// OrderRepository defines the interface for managing orders in the repository layer.
// It provides methods to retrieve, create, update, and delete orders.
type OrderRepository interface {
	// GetOrderByID retrieves an order by its ID.
	//
	// Parameters:
	//   - id: The unique identifier of the order to retrieve.
	//
	// Returns:
	//   - A pointer to the Order entity if found.
	//   - An error if the retrieval process fails or the order is not found.
	GetOrderByID(ctx context.Context, id int64) (*entity.Order, error)

	// CreateOrder creates a new order in the repository.
	//
	// Parameters:
	//   - order: A pointer to the Order entity to be created.
	//
	// Returns:
	//   - A pointer to the created Order entity with updated fields.
	//   - An error if the creation process fails.
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)

	// UpdateOrder updates an existing order in the repository.
	//
	// Parameters:
	//   - order: A pointer to the Order entity to be updated.
	//
	// Returns:
	//   - A pointer to the updated Order entity.
	//   - An error if the update process fails.
	UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)

	// DeleteOrder deletes an order by its ID from the repository.
	//
	// Parameters:
	//   - id: The unique identifier of the order to delete.
	//
	// Returns:
	//   - An error if the deletion process fails or the order is not found.
	DeleteOrder(ctx context.Context, id int64) error
}

// orderRepository is a concrete implementation of the OrderRepository interface.
// It uses an in-memory map to simulate order storage.
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates and returns a new instance of orderRepository.
//
// Returns:
//   - An instance of OrderRepository.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// GetOrderByID retrieves an order by its ID from the in-memory storage.
//
// Parameters:
//   - id: The unique identifier of the order to retrieve.
//
// Returns:
//   - A pointer to the Order entity if found.
//   - An error if the order is not found.
func (r *orderRepository) GetOrderByID(ctx context.Context, id int64) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Table("orders").WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Info().Int64("orderID", id).Msg("Order not found")
			return nil, nil
		}
		log.Logger.Error().Err(err).Int64("orderID", id).Msg("Failed to get order by ID")
		return nil, err
	}

	return &order, nil
}

// CreateOrder creates a new order in the in-memory storage.
//
// Parameters:
//   - order: A pointer to the Order entity to be created.
//
// Returns:
//   - A pointer to the created Order entity with an auto-generated ID.
//   - An error if the creation process fails.
func (r *orderRepository) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	err := r.db.Table("orders").WithContext(ctx).Create(order).Error
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create order")
		return nil, err
	}

	return order, nil
}

// UpdateOrder updates an existing order in the in-memory storage.
//
// Parameters:
//   - order: A pointer to the Order entity to be updated.
//
// Returns:
//   - A pointer to the updated Order entity.
//   - An error if the update process fails.
func (r *orderRepository) UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	err := r.db.Table("orders").WithContext(ctx).Save(order).Error
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", order.ID).Msg("Failed to update order")
		return nil, err
	}
	return order, nil
}

// DeleteOrder deletes an order by its ID from the in-memory storage.
//
// Parameters:
//   - id: The unique identifier of the order to delete.
//
// Returns:
//   - An error if the order is not found or the deletion process fails.
func (r *orderRepository) DeleteOrder(ctx context.Context, id int64) error {
	order, err := r.GetOrderByID(ctx, id)
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", id).Msg("Failed to retrieve order before deletion")
		return err
	}

	if order == nil {
		log.Logger.Warn().Int64("orderID", id).Msg("Order not found for deletion")
		return gorm.ErrRecordNotFound
	}

	err = r.db.Table("orders").WithContext(ctx).Delete(&entity.Order{}, id).Error
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", id).Msg("Failed to delete order")
		return err
	}

	return nil
}
