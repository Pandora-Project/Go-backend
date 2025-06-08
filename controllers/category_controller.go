package controllers

import (
	"net/http"
	"strconv"

	"echo-gorm-project/database"
	"echo-gorm-project/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const CATEGORIES = "/categories"
const CATEGORIES_ID = "/categories/:id"
const CAT_ERROR_NOT_FOUND = "Category not found"


func RegisterCategoryRoutes(e *echo.Echo) {
	e.POST(CATEGORIES, CreateCategory)
	e.GET(CATEGORIES, GetCategories)
	e.GET(CATEGORIES_ID, GetCategory)
	e.PUT(CATEGORIES_ID, UpdateCategory)
	e.DELETE(CATEGORIES_ID, DeleteCategory)
}

func CreateCategory(c echo.Context) error {
	category := new(models.Category)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	database.DB.Create(category)
	return c.JSON(http.StatusCreated, category)
}

// GetCategories retrieves categories with optional filtering scopes
func GetCategories(c echo.Context) error {
	// Base query with preloading products
	query := database.DB.Model(&models.Category{}).Preload("Products")

	if term := c.QueryParam("search"); term != "" {
		query = query.Scopes(models.SearchCategoryByName(term))
	}

	if minParam := c.QueryParam("min_products"); minParam != "" {
		if maxParam := c.QueryParam("max_products"); maxParam != "" {
			min, err1 := strconv.Atoi(minParam)
			max, err2 := strconv.Atoi(maxParam)
			if err1 == nil && err2 == nil {
				query = query.Scopes(models.FilterByProductCount(min, max))
			}
		}
	}

	var categories []models.Category
	if err := query.Find(&categories).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No categories found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch categories"})
	}
	return c.JSON(http.StatusOK, categories)
}

func GetCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.Preload("Products").First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": CAT_ERROR_NOT_FOUND})
	}
	return c.JSON(http.StatusOK, category)
}

func UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": CAT_ERROR_NOT_FOUND})
	}
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	database.DB.Save(&category)
	return c.JSON(http.StatusOK, category)
}

func DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": CAT_ERROR_NOT_FOUND})
	}
	database.DB.Delete(&category)
	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted"})
}
