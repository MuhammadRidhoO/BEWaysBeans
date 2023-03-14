package models

type Transaction struct {
	Id             string                           `json:"id" gorm:"type: varchar(255);PRIMARY_KEY"`
	Total          int                              `json:"total" gorm:"type: int"`
	OrderDate      string                        `json:"order_date"`
	Status_Payment string                           `json:"status_payment" gorm:"type: varchar(255)"`
	User_Id        int                              `json:"user_id" gorm:"type: int"`
	User           User                             `json:"user"`
	Order          []Order_Response_For_Transaction `json:"products" gorm:"foreignKey:Transaction_Id"`
}

type TransactionResponse struct {
	Id             string                           `json:"id" gorm:"type: varchar(255);PRIMARY_KEY"`
	OrderDate      string                        `json:"order_date"`
	Total          int                              `json:"total" gorm:"type: int"`
	Status_Payment string                           `json:"status_payment" gorm:"type: varchar(255)"`
	User_Id        int                              `json:"user_id" gorm:"type: int"`
	User           User                             `json:"user"`
	Order          []Order_Response_For_Transaction `json:"products" gorm:"foreignKey:Transaction_Id"`
}

func (TransactionResponse) TableName() string {
	return "transactions"
}
