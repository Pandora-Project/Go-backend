package controllers

import (
	"net/http"
	"strconv"

	"echo-gorm-project/database"
	"echo-gorm-project/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const PRODUCTS_ID = "/products/:id"
const PROD_NOT_FOUND = "Product not found"


func RegisterProductRoutes(e *echo.Echo) {
	e.POST("/products", CreateProduct)
	e.GET("/products", GetProducts)
	e.GET(PRODUCTS_ID, GetProduct)
	e.PUT(PRODUCTS_ID, UpdateProduct)
	e.DELETE(PRODUCTS_ID, DeleteProduct)
}

func CreateProduct(c echo.Context) error {
	product := new(models.Product)
	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	database.DB.Create(product)
	return c.JSON(http.StatusCreated, product)
}

func GetProducts(c echo.Context) error {
	query := database.DB.Model(&models.Product{}).Preload("Category")

	if minParam := c.QueryParam("min_price"); minParam != "" {
		if maxParam := c.QueryParam("max_price"); maxParam != "" {
			min, err1 := strconv.ParseFloat(minParam, 64)
			max, err2 := strconv.ParseFloat(maxParam, 64)
			if err1 == nil && err2 == nil {
				query = query.Scopes(models.FilterByPrice(min, max))
			}
		}
	}

	if term := c.QueryParam("search"); term != "" {
		query = query.Scopes(models.SearchByName(term))
	}

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No products found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products"})
	}
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.Preload("Category").First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": PROD_NOT_FOUND})
	}
	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": PROD_NOT_FOUND})
	}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	database.DB.Save(&product)
	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": PROD_NOT_FOUND})
	}
	database.DB.Delete(&product)
	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted"})
}
