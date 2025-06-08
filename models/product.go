package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Description string
	Price       float64
	CategoryID  uint
	Category    Category `gorm:"foreignKey:CategoryID"`
}

func FilterByPrice(min, max float64) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("price BETWEEN ? AND ?", min, max)
    }
}

func SearchByName(term string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("name LIKE ?", "%" + term + "%")
    }
}
