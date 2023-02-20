package handlers

import (
	"BEWaysBeans/dto"
	"BEWaysBeans/models"
	"BEWaysBeans/repositories"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type handlerProduct struct {
	ProductRepository repositories.ProductRepository
}

func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
	return &handlerProduct{ProductRepository}
}

// for view all data
func (h *handlerProduct) FindProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	product, err := h.ProductRepository.FindProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: (product)}
	json.NewEncoder(w).Encode(response)
	// fmt.Println(products)
}

func (h *handlerProduct) FilterProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	name_product := queryParams.Get("name_product")
	price, _ := strconv.Atoi(queryParams.Get("price"))
	stock, _ := strconv.Atoi(queryParams.Get("stock"))

	products, err := h.ProductRepository.FilterProducts(name_product, price, stock)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: products}
	json.NewEncoder(w).Encode(response)
	fmt.Println(products)
}

func (h *handlerProduct) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//params
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	//get data
	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	//to view success get data
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: product}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")
	// request
	price, _ := strconv.Atoi(r.FormValue("price"))
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	request := dto.Request_Product{
		Name_Product:  r.FormValue("name_product"),
		Image_Product: r.FormValue("image_product"),
		Price:         price,
		Stock:         stock,
		Descraption:   r.FormValue("descraption"),
	}

	// validation
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// image
	dataContex := r.Context().Value("Error")
	var filename string
	if dataContex == nil {
		// image
		dataContex := r.Context().Value("dataFile")
		filename = dataContex.(string)
	}
	// dataContex := r.Context().Value("dataFile") // add this code
	// filename := dataContex.(string)             // add this code
	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filename, uploader.UploadParams{Folder: "WaysBeans"})

	if err != nil {
		fmt.Println(err.Error())
	}
	product := models.Product{
		Name_Product:  request.Name_Product,
		Descraption:   request.Descraption,
		Price:         request.Price,
		Stock:         request.Stock,
		Image_Product: resp.SecureURL,
	}

	// store data
	data, err := h.ProductRepository.CreateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// get data
	productGet, err := h.ProductRepository.GetProduct(data.Id)
	productGet.Image_Product = os.Getenv("PATH_FILE") + productGet.Image_Product
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	// success
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertProductResponse(productGet)}
	json.NewEncoder(w).Encode(response)
}

// Update data
func (h *handlerProduct) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// params
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// get data
	product, err := h.ProductRepository.GetProduct(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	// validation
	dataContex := r.Context().Value("dataFile")
	var filepath string
	if dataContex.(string) != "" {
		filepath = dataContex.(string)
	}

	// create empty context
	var ctx = context.Background()

	// setup cloudinary credentials
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// create new instance of cloudinary object using cloudinary credentials
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	fmt.Println("Test cld", cld)
	// Upload file to Cloudinary
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "Hausy"})
	if err != nil {
		fmt.Println(err.Error())
	}

	// image
	product.Image_Product = resp.SecureURL

	if r.FormValue("name_product") != "" {
		product.Name_Product = r.FormValue("name_product")
	}
	Price, _ := strconv.Atoi(r.FormValue("price"))
	if Price != 0 {
		product.Price = Price
	}
	Stock, _ := strconv.Atoi(r.FormValue("stock"))
	if Stock != 0 {
		product.Stock = Stock
	}
	if r.FormValue("descraption") != "" {
		product.Descraption = r.FormValue("descraption")
	}
	// date_tripInput, _ := time.Parse("2006-01-02", r.FormValue("date_trip"))
	// if date_tripInput.IsZero() {
	// 	date_trip := trip.Date_Trip
	// 	trip.Date_Trip = date_trip
	// }

	// update data
	data, err := h.ProductRepository.UpdateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// get data
	productInserted, err := h.ProductRepository.GetProduct(data.Id)
	productInserted.Image_Product = os.Getenv("PATH_FILE") + productInserted.Image_Product
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	// success
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertProductResponse(productInserted)}
	json.NewEncoder(w).Encode(response)
}

func convertProductResponse(r models.Product) dto.Response_Product {
	return dto.Response_Product{
		Id:            r.Id,
		Name_Product:  r.Name_Product,
		Price:         r.Price,
		Stock:         r.Stock,
		Descraption:   r.Descraption,
		Image_Product: r.Image_Product,
	}
}
