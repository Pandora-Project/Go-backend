package controllers

import (
	"echo-gorm-project/database"
	"echo-gorm-project/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RegisterOrderRoutes registers endpoints for orders
func RegisterOrderRoutes(e *echo.Echo) {
	e.POST("/orders", CreateOrder)
	e.GET("/orders", GetOrders)
	e.GET("/orders/:id", GetOrder)
	e.PUT("/orders/:id", UpdateOrder)
	e.DELETE("/orders/:id", DeleteOrder)
}

// CreateOrder handles creating a new order with items
func CreateOrder(c echo.Context) error {
	var input struct {
		UserID uint `json:"user_id"`
		Items  []struct {
			ProductID uint    `json:"product_id"`
			Quantity  int     `json:"quantity"`
			UnitPrice float64 `json:"unit_price"`
		} `json:"items"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Validate items and compute total
	total := 0.0
	var orderItems []models.OrderItem
	for _, it := range input.Items {
		if it.ProductID == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "product_id must be provided"})
		}
		if it.Quantity <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "quantity must be greater than 0"})
		}

		// Check product existence
		var prod models.Product
		if err := database.DB.First(&prod, it.ProductID).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		}

		total += float64(it.Quantity) * it.UnitPrice
		orderItems = append(orderItems, models.OrderItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
		})
	}

	order := models.Order{
		UserID:   input.UserID,
		Status:   "pending",
		Total:    total,
		Items:    orderItems,
		PlacedAt: time.Now(),
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create order"})
	}
	// Preload items and products
	database.DB.Preload("Items.Product").First(&order, order.ID)
	return c.JSON(http.StatusCreated, order)
}

// GetOrders retrieves orders with optional filters via scopes
func GetOrders(c echo.Context) error {
	query := database.DB.Model(&models.Order{}).Preload("Items.Product")

	// Filter by status: ?status=paid
	if status := c.QueryParam("status"); status != "" {
		query = query.Scopes(models.FilterByStatus(status))
	}
	// Filter by user_id: ?user_id=42
	if uid := c.QueryParam("user_id"); uid != "" {
		if id, err := strconv.ParseUint(uid, 10, 64); err == nil {
			query = query.Scopes(models.FilterByUserID(uint(id)))
		}
	}

	var orders []models.Order
	if err := query.Find(&orders).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No orders found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders"})
	}
	return c.JSON(http.StatusOK, orders)
}

// GetOrder retrieves a single order by ID
func GetOrder(c echo.Context) error {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("Items.Product").First(&order, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}
	return c.JSON(http.StatusOK, order)
}

// UpdateOrder updates order status (and optionally items)
func UpdateOrder(c echo.Context) error {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}
	// Bind only status field
	var input struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	order.Status = input.Status
	database.DB.Save(&order)
	database.DB.Preload("Items.Product").First(&order, order.ID)
	return c.JSON(http.StatusOK, order)
}

// DeleteOrder deletes an order by its ID using a scope
func DeleteOrder(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID"})
	}
	// Delete via scope
	result := database.DB.
		Scopes(models.FilterOrderByID(uint(id))).
		Delete(&models.Order{})
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete order"})
	}
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Order deleted"})
}
