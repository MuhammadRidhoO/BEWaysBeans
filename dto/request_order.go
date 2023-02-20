package dto

import "BEWaysBeans/models"

type Request_Order struct {
	Product_Id int `json:"product_id"`
}
type Update_Request_Order struct {
	Qty   int    `json:"qty"`
	Event string `json:"event"`
}
type Order_Response struct {
	Id      int            `json:"id"`
	Qty     int            `json:"qty"`
	Product models.Product `json:"product"`
	User_Id int            `json:"user_id"`
}
