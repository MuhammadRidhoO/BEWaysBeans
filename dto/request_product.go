package dto

type Request_Product struct {
	Name_Product  string `json:"name_product" form:"name_product"`
	Price         int    `json:"price" form:"price"`
	Descraption   string `json:"descraption" form:"descraption"`
	Stock         int    `json:"stock" form:"stock"`
	Image_Product string `json:"image_product" form:"image_product"`
}
type Request_Update_Product struct {
	Name_Product  string `json:"name_product" form:"name_product"`
	Price         int    `json:"price" form:"price"`
	Descraption   string `json:"descraption" form:"descraption"`
	Stock         int    `json:"stock" form:"stock"`
	Image_Product string `json:"image_product" form:"image_product"`
}
