package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID       uint           `gorm:"primaryKey"`
	Products []Product      `gorm:"many2many:cart_products;"`
	ProductID uint
	Quantity  int
	Product   Product       `gorm:"foreignKey:ProductID"`
}

type CartProduct struct {
    gorm.Model
    ProductID uint `json:"product_id"`
    Quantity  int  `json:"quantity"`
}

func FilterByQuantity(minQuantity, maxQuantity int) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("quantity BETWEEN ? AND ?", minQuantity, maxQuantity)
    }
}

func FilterByProductID(productID uint) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("product_id = ?", productID)
    }
}

func FilterCartByID(id uint) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("id = ?", id)
    }
}