package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
    gorm.Model
    UserID     uint         `json:"user_id"`      
    Status     string       `json:"status"`       
    Total      float64      `json:"total"`        
    Items      []OrderItem  `gorm:"foreignKey:OrderID" json:"items"`
    PlacedAt   time.Time    `json:"placed_at"`
}

func FilterByStatus(status string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("status = ?", status)
    }
}

func FilterByUserID(userID uint) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("user_id = ?", userID)
    }
}

func FilterOrderByID(orderID uint) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("id = ?", orderID)
    }
}