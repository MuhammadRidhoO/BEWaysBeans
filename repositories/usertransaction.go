package repositories

import (
	"BEWaysBeans/models"

	"gorm.io/gorm"
)

type UserTrcRepository interface {
	FindUserTrc(ID int) ([]models.Transaction, error)
}

func RepositoryUserTrc(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindUserTrc(ID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Order").Preload("Order.Product").Where("user_id = ?", ID).Order("order_date desc").Find(&transactions).Error

	return transactions, err
}