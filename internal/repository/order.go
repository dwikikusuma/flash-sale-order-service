package repository

import "order-service/internal/entity"

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
	GetOrderByID(id int64) (*entity.Order, error)

	// CreateOrder creates a new order in the repository.
	//
	// Parameters:
	//   - order: A pointer to the Order entity to be created.
	//
	// Returns:
	//   - A pointer to the created Order entity with updated fields.
	//   - An error if the creation process fails.
	CreateOrder(order *entity.Order) (*entity.Order, error)

	// UpdateOrder updates an existing order in the repository.
	//
	// Parameters:
	//   - order: A pointer to the Order entity to be updated.
	//
	// Returns:
	//   - A pointer to the updated Order entity.
	//   - An error if the update process fails.
	UpdateOrder(order *entity.Order) (*entity.Order, error)

	// DeleteOrder deletes an order by its ID from the repository.
	//
	// Parameters:
	//   - id: The unique identifier of the order to delete.
	//
	// Returns:
	//   - An error if the deletion process fails or the order is not found.
	DeleteOrder(id int64) error
}

// orderRepository is a concrete implementation of the OrderRepository interface.
// It uses an in-memory map to simulate order storage.
type orderRepository struct {
}

// NewOrderRepository creates and returns a new instance of orderRepository.
//
// Returns:
//   - An instance of OrderRepository.
func NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

// orders is an in-memory map simulating order storage.
// In a real application, this would be replaced with a database or other persistent storage.
var orders = map[int64]*entity.Order{
	1: {
		ID:              1,
		UserID:          123,
		ProductRequests: make([]entity.OrderRequest, 0),
		Quantity:        2,
		TotalPrice:      100.0,
		Status:          "created",
	},
	2: {
		ID:              2,
		UserID:          456,
		ProductRequests: make([]entity.OrderRequest, 0),
		Quantity:        1,
		TotalPrice:      50.0,
		Status:          "created",
	},
}

// GetOrderByID retrieves an order by its ID from the in-memory storage.
//
// Parameters:
//   - id: The unique identifier of the order to retrieve.
//
// Returns:
//   - A pointer to the Order entity if found.
//   - An error if the order is not found.
func (r *orderRepository) GetOrderByID(id int64) (*entity.Order, error) {
	order, ok := orders[id]
	if !ok {
		return nil, nil
	}
	return order, nil
}

// CreateOrder creates a new order in the in-memory storage.
//
// Parameters:
//   - order: A pointer to the Order entity to be created.
//
// Returns:
//   - A pointer to the created Order entity with an auto-generated ID.
//   - An error if the creation process fails.
func (r *orderRepository) CreateOrder(order *entity.Order) (*entity.Order, error) {
	order.ID = int64(len(orders) + 1) // Simulating an auto-generated ID
	orders[order.ID] = order
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
func (r *orderRepository) UpdateOrder(order *entity.Order) (*entity.Order, error) {
	orders[order.ID] = order
	return order, nil
}

// DeleteOrder deletes an order by its ID from the in-memory storage.
//
// Parameters:
//   - id: The unique identifier of the order to delete.
//
// Returns:
//   - An error if the order is not found or the deletion process fails.
func (r *orderRepository) DeleteOrder(id int64) error {
	_, ok := orders[id]
	if !ok {
		return nil // Order not found, nothing to delete
	}
	delete(orders, id)
	return nil
}
