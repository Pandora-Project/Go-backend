package controllers

import (
	"bytes"
	"echo-gorm-project/database"
	"echo-gorm-project/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCategoryTestDB() {
    db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    database.DB = db
    database.DB.AutoMigrate(&models.Category{}, &models.Product{})
}

func setupCategoryEcho() *echo.Echo {
    e := echo.New()
    RegisterCategoryRoutes(e)
    return e
}

func TestCreateCategory_Success(t *testing.T) {
    setupCategoryTestDB()
    e := setupCategoryEcho()

    body := map[string]interface{}{
        "name": "TestCategory",
    }
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if err := CreateCategory(c); err != nil {
        t.Fatalf("CreateCategory failed: %v", err)
    }
    if rec.Code != http.StatusCreated {
        t.Errorf("expected 201, got %d", rec.Code)
    }
}

func TestCreateCategory_BadRequest(t *testing.T) {
    setupCategoryTestDB()
    e := setupCategoryEcho()

    req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("bad json")))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    CreateCategory(c)
    if rec.Code != http.StatusBadRequest {
        t.Errorf("expected 400, got %d", rec.Code)
    }
}

func TestGetCategory_NotFound(t *testing.T) {
    setupCategoryTestDB()
    e := setupCategoryEcho()

    req := httptest.NewRequest(http.MethodGet, "/categories/999", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("999")

    GetCategory(c)
    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", rec.Code)
    }
}

func TestUpdateCategory_NotFound(t *testing.T) {
    setupCategoryTestDB()
    e := setupCategoryEcho()

    body := map[string]interface{}{
        "name": "UpdatedCategory",
    }
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPut, "/categories/999", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("999")

    UpdateCategory(c)
    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", rec.Code)
    }
}

func TestDeleteCategory_NotFound(t *testing.T) {
    setupCategoryTestDB()
    e := setupCategoryEcho()

    req := httptest.NewRequest(http.MethodDelete, "/categories/999", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("999")

    DeleteCategory(c)
    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404, got %d", rec.Code)
    }
}

func TestGetCategories_Empty(t *testing.T) {
    setupCategoryTestDB()
    e := setupCategoryEcho()

    req := httptest.NewRequest(http.MethodGet, "/categories", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    GetCategories(c)
    if rec.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rec.Code)
    }
}