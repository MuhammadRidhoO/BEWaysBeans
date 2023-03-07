package dto

import "BEWaysBeans/models"

type ProductRequestForTransaction struct {
	Id         int `json:"id"`
	Product_Id int `json:"product_id"`
	Qty        int `json:"qty"`
}

type ProductResponseForTransaction struct {
	Id            int    `json:"id"`
	Name_Product  string `json:"name_product"`
	Price         int    `json:"price"`
	Image_Product string `json:"image_product"`
	Descraption   string `json:"descraption"`
	Qty           int    `json:"qty"`
}

type CreateTransactionRequest struct {
	Id       int                            `json:"id"`
	Total    int                            `json:"total" validate:"required"`
	User_Id  int                            `json:"user_id" validate:"required"`
	Products []ProductRequestForTransaction `json:"products" validate:"required"`
}
type UpdateTransactionRequest struct {
	Total          int                            `json:"total"`
	User_Id        int                            `json:"user_id"`
	Products       []ProductRequestForTransaction `json:"products"`
	Status_Payment string                         `json:"status_payment"`
}

type TransactionResponse struct {
	Id             int                             `json:"id"`
	Total          int                             `json:"total"`
	Status_Payment string                          `json:"status_payment"`
	User           models.User                     `json:"user"`
	Products       []ProductResponseForTransaction `json:"products"`
}
