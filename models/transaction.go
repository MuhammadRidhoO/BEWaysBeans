package models

import "time"

type Transaction struct {
	Id             int                           `json:"id" gorm:"type: int;PRIMARY_KEY"`
	Total          int                           `json:"total" gorm:"type: int"`
	Status_Payment string                        `json:"status_payment" gorm:"type: varchar(255)"`
	User_Id        int                           `json:"user_id" gorm:"type: int"`
	User           User                          `json:"user"`
	Order          []Order_Response_For_Transaction `json:"products" gorm:"foreignKey:Transaction_Id"`
}

type TransactionResponse struct {
	Id             int                           `json:"id" gorm:"type: int;PRIMARY_KEY"`
	OrderDate      time.Time                     `json:"order_date"`
	Total          int                           `json:"total" gorm:"type: int"`
	Status_Payment string                        `json:"status_payment" gorm:"type: varchar(255)"`
	User_Id        int                           `json:"user_id" gorm:"type: int"`
	User           User                          `json:"user"`
	Order          []Order_Response_For_Transaction `json:"products" gorm:"foreignKey:Transaction_Id"`
}

func (TransactionResponse) TableName() string {
	return "transactions"
}