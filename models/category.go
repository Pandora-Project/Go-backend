package models

import "gorm.io/gorm"

type Category struct {
    ID       uint      `json:"id" gorm:"primaryKey"`
    Name     string    `json:"name"`
    Products []Product `json:"products" gorm:"foreignKey:CategoryID"`
}

func SearchCategoryByName(term string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("name LIKE ?", "%"+term+"%")
    }
}

func FilterByProductCount(min, max int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("(SELECT COUNT(*) FROM products WHERE products.category_id = categories.id) BETWEEN ? AND ?", min, max)
    }
}