package models

import "gorm.io/gorm"

type OrderItem struct {
    gorm.Model
    OrderID   uint     `json:"order_id"`
    ProductID uint     `json:"product_id"`
    Quantity  int      `json:"quantity"`
    UnitPrice float64  `json:"unit_price"`  // snapshot of product price at purchase
    Product   Product  `gorm:"foreignKey:ProductID" json:"product"`
}
