// main.go

package main

import (
	"echo-gorm-project/controllers"
	"echo-gorm-project/database"
	"echo-gorm-project/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    // Init DB
    database.ConnectDatabase()
    database.DB.AutoMigrate(
        &models.Product{},
        &models.Category{},
        &models.Cart{},
        &models.Order{},      
        &models.OrderItem{},
    )

    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Register your existing routes
    controllers.RegisterProductRoutes(e)
    controllers.RegisterCategoryRoutes(e)
    controllers.RegisterCartRoutes(e)

    // **Register the orders routes**
    controllers.RegisterOrderRoutes(e)

    e.Logger.Fatal(e.Start(":8080"))
}
