package controllers

import (
	"bytes"
	"echo-gorm-project/database"
	"echo-gorm-project/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupOrderTestDB() {
    db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    database.DB = db
    database.DB.AutoMigrate(&models.Order{}, &models.OrderItem{}, &models.Product{})
}

func setupOrderEcho() *echo.Echo {
    e := echo.New()
    RegisterOrderRoutes(e)
    return e
}

func TestCreateOrder_BadRequest(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte("bad json")))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    CreateOrder(c)
    if rec.Code != http.StatusBadRequest {
        t.Errorf("expected 400, got %d", rec.Code)
    }
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    body := map[string]interface{}{
        "user_id": 1,
        "items": []map[string]interface{}{
            {"product_id": 999, "quantity": 1, "unit_price": 10},
        },
    }
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    CreateOrder(c)
    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", rec.Code)
    }
}

func TestCreateOrder_Success(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    // Add a product to reference
    product := models.Product{Name: "TestProduct", Price: 10}
    database.DB.Create(&product)

    body := map[string]interface{}{
        "user_id": 1,
        "items": []map[string]interface{}{
            {"product_id": product.ID, "quantity": 2, "unit_price": 10},
        },
    }
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    CreateOrder(c)
    if rec.Code != http.StatusCreated {
        t.Errorf("expected 201, got %d", rec.Code)
    }
}

func TestGetOrder_NotFound(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    req := httptest.NewRequest(http.MethodGet, "/orders/999", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("999")

    GetOrder(c)
    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", rec.Code)
    }
}

func TestGetOrder_Success(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    // Create order and product
    product := models.Product{Name: "TestProduct", Price: 10}
    database.DB.Create(&product)
    order := models.Order{
        UserID:   1,
        Status:   "pending",
        Total:    20,
        PlacedAt: time.Now(),
        Items: []models.OrderItem{
            {ProductID: product.ID, Quantity: 2, UnitPrice: 10},
        },
    }
    database.DB.Create(&order)

    req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("1")

    GetOrder(c)
    if rec.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rec.Code)
    }
}

func TestUpdateOrder_NotFound(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    body := map[string]interface{}{
        "status": "paid",
    }
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPut, "/orders/999", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("999")

    UpdateOrder(c)
    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", rec.Code)
    }
}

func TestUpdateOrder_Success(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    // Create order and product
    product := models.Product{Name: "TestProduct", Price: 10}
    database.DB.Create(&product)
    order := models.Order{
        UserID:   1,
        Status:   "pending",
        Total:    20,
        PlacedAt: time.Now(),
        Items: []models.OrderItem{
            {ProductID: product.ID, Quantity: 2, UnitPrice: 10},
        },
    }
    database.DB.Create(&order)

    body := map[string]interface{}{
        "status": "paid",
    }
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPut, "/orders/1", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("1")

    UpdateOrder(c)
    if rec.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rec.Code)
    }
}

func TestDeleteOrder_NotFound(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    req := httptest.NewRequest(http.MethodDelete, "/orders/999", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("999")

    DeleteOrder(c)
    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", rec.Code)
    }
}

func TestDeleteOrder_Success(t *testing.T) {
    setupOrderTestDB()
    e := setupOrderEcho()

    // Create order and product
    product := models.Product{Name: "TestProduct", Price: 10}
    database.DB.Create(&product)
    order := models.Order{
        UserID:   1,
        Status:   "pending",
        Total:    20,
        PlacedAt: time.Now(),
        Items: []models.OrderItem{
            {ProductID: product.ID, Quantity: 2, UnitPrice: 10},
        },
    }
    database.DB.Create(&order)

    req := httptest.NewRequest(http.MethodDelete, "/orders/1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("1")

    DeleteOrder(c)
    if rec.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rec.Code)
    }
}