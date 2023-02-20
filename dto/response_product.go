package dto

type Response_Product struct {
	Id            int    `json:"id" form:"id"`
	Name_Product  string `json:"name_product" form:"name_product"`
	Price         int    `json:"price" form:"price"`
	Descraption   string `json:"descraption" form:"descraption"`
	Stock         int    `json:"stock" form:"stock"`
	Image_Product string `json:"image_product" form:"image_product"`
}
type Response_Update_Product struct {
	Id            int    `json:"id" form:"id"`
	Name_Product  string `json:"name_product" form:"name_product"`
	Price         int    `json:"price" form:"price"`
	Descraption   string `json:"descraption" form:"descraption"`
	Stock         int    `json:"stock" form:"stock"`
	Image_Product string `json:"image_product" form:"image_product"`
}
