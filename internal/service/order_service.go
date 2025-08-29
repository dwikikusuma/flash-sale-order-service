package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/infrastructure/log"
	"order-service/internal/entity"
	"order-service/internal/repository"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type OrderService interface {
	// CreateOrder creates a new order with an initial status of "created".
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	// UpdateOrder updates an existing order by modifying its status to "updated".
	UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	// CancelOrder cancels an existing order by modifying its status to "cancelled".
	CancelOrder(ctx context.Context, orderId int64) (*entity.Order, error)
}

// orderService provides methods to manage orders, including creating, updating, and canceling orders.
type orderService struct {
	OrderRepository   repository.OrderRepository
	ProductServiceURL string // URL for the product service, if needed for communication
	PricingServiceURL string // URL for the pricing service, if needed for communication
	KafkaWriter       *kafka.Writer
}

// NewOrderService creates and returns a new instance of orderService.
func NewOrderService(productRepository repository.OrderRepository, productServiceURL, PricingServiceURL string, kafkaWriter *kafka.Writer) OrderService {
	return &orderService{
		OrderRepository:   productRepository,
		ProductServiceURL: productServiceURL,
		PricingServiceURL: PricingServiceURL,
		KafkaWriter:       kafkaWriter,
	}
}

// CreateOrder creates a new order with an initial status of "created".
// It simulates assigning an auto-generated ID to the order.
//
// Parameters:
//   - order: A pointer to the Order entity to be created.
//
// Returns:
//   - A pointer to the created Order entity with updated fields.
//   - An error if the creation process fails.
func (s *orderService) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	// Logic to create an order
	// This could involve saving the order to a database, etc.
	var totalPrice float64

	availabilityCh := make(chan entity.AvailabilityChannel, len(order.ProductRequests))
	pricingCh := make(chan entity.PricingChannel, len(order.ProductRequests))

	// Launch goroutines to fetch availability and pricing data concurrently
	for _, productRequest := range order.ProductRequests {
		go func(productRequest *entity.OrderRequest) {
			available, err := s.checkProductStock(productRequest.ProductID, productRequest.Quantity)
			availabilityCh <- entity.AvailabilityChannel{
				ProductID: productRequest.ProductID,
				Available: available,
				Error:     err,
			}
		}(&productRequest)

		go func(productRequest *entity.OrderRequest) {
			pricing, err := s.getPricing(productRequest.ProductID)
			pricingCh <- entity.PricingChannel{
				ProductID:  productRequest.ProductID,
				FinalPrice: pricing.FinalPrice,
				MarkUp:     pricing.MarkUp,
				Discount:   pricing.Discount,
				Error:      err,
			}
		}(&productRequest)
	}

	// Process results from channels
	// NOTE: Current design doesn't require mapping between availability and pricing channels
	// since we process them independently (availability for validation, pricing by ProductID matching).
	// However, for future development where you need to correlate results from multiple channels
	// for the same product, consider using a map-based approach or combined result channels
	// to ensure proper pairing of related data.
	for range order.ProductRequests {
		availabilityResult := <-availabilityCh
		pricingResult := <-pricingCh

		if availabilityResult.Error != nil {
			log.Logger.Error().Err(availabilityResult.Error).Int64("productID", availabilityResult.ProductID).Msg("Failed to check product stock")
			return nil, fmt.Errorf("failed to check product stock for product ID %d: %w", availabilityResult.ProductID, availabilityResult.Error)
		}
		if !availabilityResult.Available {
			log.Logger.Warn().Int64("productID", availabilityResult.ProductID).Msg("Insufficient stock for product")
			return nil, fmt.Errorf("insufficient stock for product ID %d", availabilityResult.ProductID)
		}
		if pricingResult.Error != nil {
			log.Logger.Error().Err(pricingResult.Error).Int64("productID", pricingResult.ProductID).Msg("Failed to get pricing for product")
			return nil, fmt.Errorf("failed to get pricing for product ID %d: %w", pricingResult.ProductID, pricingResult.Error)
		}

		for _, productRequest := range order.ProductRequests {
			if productRequest.ProductID == pricingResult.ProductID {
				productRequest.Discount = pricingResult.Discount
				productRequest.MarkUp = pricingResult.MarkUp
				productRequest.FinalPrice = pricingResult.FinalPrice
				totalPrice += productRequest.FinalPrice
			}
		}
	}

	err := s.OrderRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		err := s.OrderRepository.CreateOrderTx(ctx, tx, order)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to create order in transaction")
			return fmt.Errorf("failed to create order in transaction: %w", err)
		}

		orderRequests := s.mapOrderRequestWithOrderID(order)
		err = s.OrderRepository.CreateOrderRequestTx(ctx, tx, orderRequests)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to create order requests in transaction")
			return fmt.Errorf("failed to create order requests in transaction: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Logger.Error().Err(err).Msg("Transaction failed, rolling back")
		return nil, err
	}

	err = s.publishOrderCreatedEvent(order, "created")
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", order.ID).Msg("Failed to publish order created event")
		return nil, fmt.Errorf("failed to publish order created event: %w", err)
	}

	return order, nil
}

// UpdateOrder updates an existing order by modifying its status to "updated".
//
// Parameters:
//   - order: A pointer to the Order entity to be updated.
//
// Returns:
//   - A pointer to the updated Order entity.
//   - An error if the update process fails.
func (s *orderService) UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	// Logic to update an existing order
	// This could involve updating the order in a database, etc.

	if order.Status == "Paid" {
		for _, orderRequest := range order.ProductRequests {
			match, err := s.checkProductStock(orderRequest.ProductID, orderRequest.Quantity)
			if err != nil {
				log.Logger.Error().Err(err).Int64("productID", orderRequest.ProductID).Msg("Failed to check product stock during order update")
				return nil, fmt.Errorf("failed to check product stock for product ID %d: %w", orderRequest.ProductID, err)
			}

			if !match {
				log.Logger.Warn().Int64("productID", orderRequest.ProductID).Msg("Insufficient stock for product during order update")
				return nil, fmt.Errorf("insufficient stock for product ID %d", orderRequest.ProductID)
			}
		}
	}

	updatedOrder, err := s.OrderRepository.UpdateOrder(ctx, order)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to update order")
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	if updatedOrder == nil {
		log.Logger.Warn().Int64("orderID", order.ID).Msg("Order not found for update")
		return nil, fmt.Errorf("order with ID %d not found", order.ID)
	}

	err = s.publishOrderCreatedEvent(updatedOrder, "updated")
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", updatedOrder.ID).Msg("Failed to publish order updated event")
		return nil, fmt.Errorf("failed to publish order updated event: %w", err)
	}

	return updatedOrder, nil
}

// CancelOrder cancels an existing order by modifying its status to "cancelled".
//
// Parameters:
//   - orderId: The ID of the order to be canceled.
//
// Returns:
//   - A pointer to the canceled Order entity.
//   - An error if the cancellation process fails.
func (s *orderService) CancelOrder(ctx context.Context, orderId int64) (*entity.Order, error) {
	// Logic to cancel an order
	// This could involve updating the order status in a database, etc.
	order, err := s.OrderRepository.GetOrderByID(ctx, orderId)
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", orderId).Msg("Failed to retrieve order for cancellation")
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	if order == nil {
		log.Logger.Warn().Int64("orderID", orderId).Msg("Order not found for cancellation")
		return nil, fmt.Errorf("order with ID %d not found", orderId)
	}

	order.Status = "cancelled" // Simulating a cancellation of the order
	cancelledOrder, err := s.OrderRepository.UpdateOrder(ctx, order)
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", orderId).Msg("Failed to cancel order")
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	err = s.publishOrderCreatedEvent(cancelledOrder, "cancelled")
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", cancelledOrder.ID).Msg("Failed to publish order cancelled event")
		return nil, fmt.Errorf("failed to publish order cancelled event: %w", err)
	}

	return cancelledOrder, nil
}

func (s *orderService) checkProductStock(productID int64, quantity int64) (bool, error) {
	response, err := http.Get(fmt.Sprintf("%s/product/%d/stock", s.ProductServiceURL, productID))
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to check product stock")
		return false, fmt.Errorf("failed to check product stock: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Logger.Error().Int64("productID", productID).Int("statusCode", response.StatusCode).Msg("Failed to check product stock")
		return false, fmt.Errorf("failed to check product stock, status code: %d", response.StatusCode)
	}

	var stockResponse map[string]int
	err = json.NewDecoder(response.Body).Decode(&stockResponse)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to decode stock response")
		return false, fmt.Errorf("failed to decode stock response: %w", err)
	}

	productStock, exists := stockResponse["stock"]
	if !exists {
		log.Logger.Warn().Int64("productID", productID).Msg("Stock information not found for product")
		return false, fmt.Errorf("stock information not found for product ID %d", productID)
	}

	return productStock >= int(quantity), nil
}

func (s *orderService) getPricing(productID int64) (*entity.Pricing, error) {
	response, err := http.Get(fmt.Sprintf("%s/product/%d/price", s.PricingServiceURL, productID))
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to get product pricing")
		return nil, fmt.Errorf("failed to get product pricing: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Logger.Error().Int64("productID", productID).Int("statusCode", response.StatusCode).Msg("Failed to get product pricing")
		return nil, fmt.Errorf("failed to get product pricing, status code: %d", response.StatusCode)
	}

	var pricing entity.Pricing
	err = json.NewDecoder(response.Body).Decode(&pricing)
	if err != nil {
		log.Logger.Error().Err(err).Int64("productID", productID).Msg("Failed to decode pricing response")
		return nil, fmt.Errorf("failed to decode pricing response: %w", err)
	}

	return &pricing, nil
}

func (s *orderService) publishOrderCreatedEvent(order *entity.Order, key string) error {
	orderJson, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("order.%s.%d", key, order.ID)),
		Value: orderJson,
	}

	err = s.KafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Logger.Error().Err(err).Int64("orderID", order.ID).Msg("Failed to publish order created event to Kafka")
		return fmt.Errorf("failed to publish order created event to Kafka: %w", err)
	}

	return nil
}

func (s *orderService) mapOrderRequestWithOrderID(order *entity.Order) []entity.OrderRequest {
	var orderRequests []entity.OrderRequest
	for _, productRequest := range order.ProductRequests {
		productRequest.OrderID = order.ID
		orderRequests = append(orderRequests, productRequest)
	}
	return orderRequests
}
