package models

type Order struct {
	Id             int
	User_Id        int  `json:"user_id" gorm:"type: int"`
	User           User `json:"user"`
	Transaction_Id int  `gorm:"type: int"`
	Transaction    TransactionResponse
	Product_Id     int `gorm:"type: int"`
	Product        Product
	Qty            int `gorm:"type: int"`
}

type Order_Response_For_Transaction struct {
	Id             int `json:"-"`
	Transaction_Id int `json:"-" gorm:"type: int"`
	Product_Id     int `json:"-"`
	Product        Product
	Qty            int `json:"Qty" gorm:"type: int"`
}

func (Order_Response_For_Transaction) TableName() string {
	return "orders"
}
