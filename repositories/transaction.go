package repositories

import (
	"BEWaysBeans/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)
	FindTransactionsByUser(Id int) ([]models.Transaction, error)
	GetTransaction(ID string) (models.Transaction, error)
	GetTransactionString(ID string) (models.Transaction, error)
	CreateTransaction(newTransaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(Status_Payment string, ID string) (models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload("Order").Preload("Order.Product").Order("order_date desc").Find(&transaction).Error

	return transaction, err
}

// mengambil semua transaksi berdasarkan user tertentu
func (r *repository) FindTransactionsByUser(Id int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload("Order").Preload("Order.Product").Where("user_id = ?", Id).Find(&transaction).Error

	return transaction, err
}

// mengambil 1 transaksi berdasarkan id
func (r *repository) GetTransaction(ID string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Order").Preload("Order.Product").First(&transaction,"id = ?", ID).Error

	return transaction, err
}

// mengambil 1 transaksi berdasarkan id
func (r *repository) GetTransactionString(ID string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Order").Preload("Order.Product").First(&transaction,"id = ?", ID).Error

	return transaction, err
}

// menambahkan transaksi baru
func (r *repository) CreateTransaction(newTransaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&newTransaction).Error

	return newTransaction, err
}

// mengupdate status transaksi berdasarkan id
func (r *repository) UpdateTransaction(Status_Payment string, ID string) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("User").Preload("Order").Preload("Order.Product").First(&transaction, "id = ?", ID)

	// If is different & Status is "success" decrement available quota on data trip
	if Status_Payment != transaction.Status_Payment && Status_Payment == "success" {
		for _, order := range transaction.Order {
			var product models.Product
			r.db.First(&product, order.Product.Id)
			product.Stock = product.Stock - order.Qty
			r.db.Model(&product).Updates(product)
		}
	}

	// If is different & Status is "reject" decrement available quota on data trip
	if Status_Payment != transaction.Status_Payment && Status_Payment == "rejected" {
		for _, order := range transaction.Order {
			var product models.Product
			r.db.First(&product, order.Product.Id)
			product.Stock = product.Stock + order.Qty
			r.db.Model(&product).Updates(product)
		}
	}

	// change transaction status
	// transaction.Status_Payment = Status_Payment

	// fmt.Println(status)
	// fmt.Println(transaction.Status)
	// fmt.Println(transaction.ID)

	err := r.db.Model(&transaction).Where("id = ?", ID).Update("status_payment", Status_Payment).Error

	return transaction, err
}
