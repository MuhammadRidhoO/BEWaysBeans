package repositories

import (
	"BEWaysBeans/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	FindOrders(User_Id int) ([]models.Order, error)
	GetOrder(ID int) (models.Order, error)
	GetOrderByProduct(Product_Id int, User_Id int) (models.Order, error)
	Create_Order(newOrder models.Order) (models.Order, error)
	UpdateOrder(order models.Order) (models.Order, error)
	DeleteOrder(order models.Order) (models.Order, error)
	DeleteAll(order models.Order, ID int) (models.Order, error)
}

func RepositoryOrder(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindOrders(User_Id int) ([]models.Order, error) {
	var Orders []models.Order
	err := r.db.Preload("Product").Where("user_id = ?", User_Id).Where("transaction_id IS NULL").Find(&Orders).Error

	return Orders, err
}

func (r *repository) GetOrder(ID int) (models.Order, error) {
	var Order models.Order
	// not yet using category relation, cause this step doesnt Belong to Many
	err := r.db.Preload("Product").Preload("User").First(&Order, "id = ?", ID).Error

	return Order, err
}

func (r *repository) GetOrderByProduct(Product_Id int, User_Id int) (models.Order, error) {
	var order models.Order
	err := r.db.Preload("Product").Where("user_id = ?", User_Id).Where("transaction_id IS NULL").First(&order, "product_id = ?", Product_Id).Error
	return order, err
}

func (r *repository) Create_Order(newOrder models.Order) (models.Order, error) {
	err := r.db.Select("Product_Id", "Qty", "User_Id").Create(&newOrder).Error

	return newOrder, err
}

func (r *repository) UpdateOrder(order models.Order) (models.Order, error) {
	err := r.db.Model(&order).Updates(order).Error
	return order, err
}

// menghapus order
func (r *repository) DeleteOrder(order models.Order) (models.Order, error) {
	err := r.db.Delete(&order).Error
	return order, err
}

func (r *repository) DeleteAll(order models.Order, ID int) (models.Order, error) {
	err := r.db.Where("user_id = ?", ID).Delete(&order).Error
	return order, err
}
