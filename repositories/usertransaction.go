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
	err := r.db.Preload("Product").Preload("User").Where("user_id =?", ID).Find(&transactions).Error

	return transactions, err
}