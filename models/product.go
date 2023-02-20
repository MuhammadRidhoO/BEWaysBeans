package models

type Product struct {
	Id            int    `json:"id"`
	Name_Product  string `json:"name_product"`
	Price         int    `json:"price"`
	Descraption   string `json:"description"`
	Stock         int    `json:"stock"`
	Image_Product string `json:"image_product"`
}

func (Product) TableName() string {
	return "products"
}