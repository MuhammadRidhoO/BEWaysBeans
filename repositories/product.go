package repositories

import (
	"BEWaysBeans/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindProducts() ([]models.Product, error)
	GetProduct(ID int) (models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	FilterProducts(Name_Product string, price int, Stock int) ([]models.Product, error)
	UpdateProduct(Product models.Product) (models.Product, error)
	DeleteProduct(Product models.Product) (models.Product, error)
}

func RepositoryProduct(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindProducts() ([]models.Product, error) {
	var product []models.Product
	err := r.db.Find(&product).Error

	return product, err
}

func (r *repository) GetProduct(ID int) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, ID).Error

	return product, err
}
func (r *repository) CreateProduct(product models.Product) (models.Product, error) {
	err := r.db.Create(&product).Error

	return product, err
}

func (r *repository) FilterProducts(Name_Product string, price int, Stock int) ([]models.Product, error) {
	var products []models.Product
	var err error

	if Name_Product != "" {
		err = r.db.Where("name_product = ?", Name_Product).Find(&products).Error

	} else if price != 0 {

		err = r.db.Where("price <= ?", price).Find(&products).Error

	} else if Stock != 0 {
		err = r.db.Where("stock = ? ", Stock).Find(&products).Error

	} else {
		err = r.db.Find(&products).Error
	}

	return products, err
}

func (r *repository) UpdateProduct(product models.Product) (models.Product, error) {
	err := r.db.Save(&product).Error

	return product, err
}
func (r *repository) DeleteProduct(product models.Product) (models.Product, error) {
	err := r.db.Delete(&product).Error

	return product, err
}
