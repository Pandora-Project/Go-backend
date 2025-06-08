package controllers

import (
	"net/http"
	"strconv"

	"echo-gorm-project/database"
	"echo-gorm-project/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterCartRoutes(e *echo.Echo) {
	e.POST("/carts", AddCartItem)
	e.GET("/carts", GetCartItems)
	e.DELETE("/carts/:id", DeleteCartItem)
}

func AddCartItem(c echo.Context) error {
	var input struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if input.ProductID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "product_id must be provided"})
	}
	if input.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "quantity must be greater than 0"})
	}

	var product models.Product
	if err := database.DB.First(&product, input.ProductID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	cartItem := models.Cart{
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}
	if err := database.DB.Create(&cartItem).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add item to cart"})
	}

	database.DB.Preload("Product").First(&cartItem, cartItem.ID)
	return c.JSON(http.StatusCreated, cartItem)
}

func GetCartItems(c echo.Context) error {
	query := database.DB.Model(&models.Cart{}).Preload("Product")

	if pid := c.QueryParam("product_id"); pid != "" {
		if id, err := strconv.ParseUint(pid, 10, 64); err == nil {
			query = query.Scopes(models.FilterByProductID(uint(id)))
		}
	}

	minQ := c.QueryParam("min_quantity")
	maxQ := c.QueryParam("max_quantity")
	if minQ != "" && maxQ != "" {
		minVal, err1 := strconv.Atoi(minQ)
		maxVal, err2 := strconv.Atoi(maxQ)
		if err1 == nil && err2 == nil {
			query = query.Scopes(models.FilterByQuantity(minVal, maxVal))
		}
	}

	var items []models.Cart
	if err := query.Find(&items).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No cart items found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch cart items"})
	}
	return c.JSON(http.StatusOK, items)
}

func DeleteCartItem(c echo.Context) error {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cart ID"})
    }

    result := database.DB.
        Scopes(models.FilterCartByID(uint(id))).
        Delete(&models.Cart{})

    if result.Error != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete cart item"})
    }

    if result.RowsAffected == 0 {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "Cart item not found"})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Cart item deleted"})
}
