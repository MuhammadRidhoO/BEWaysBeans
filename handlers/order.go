package handlers

import (
	"BEWaysBeans/dto"
	"BEWaysBeans/models"
	"BEWaysBeans/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerOrder struct {
	OrderRepository repositories.OrderRepository
}

func HandlerOrder(orderRepository repositories.OrderRepository) *handlerOrder {
	return &handlerOrder{orderRepository}
}

func (h *handlerOrder) FindOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	UserId := int(userInfo["id"].(float64))

	order, err := h.OrderRepository.FindOrders(UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertMultipleOrderResponse(order)}
	json.NewEncoder(w).Encode(response)
	// fmt.Println(products)
}

func (h *handlerOrder) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//params
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	//get data
	order, err := h.OrderRepository.GetOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	//to view success get data
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: order}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerOrder) Create_Order_Product(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request dto.Request_Order
	json.NewDecoder(r.Body).Decode(&request)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	idUser := int(userInfo["id"].(float64))

	// periksa order dengan product id yang sama
	order, err := h.OrderRepository.GetOrderByProduct(request.Product_Id, idUser)
	if err != nil {
		// bila belum ada, maka buat baru
		newOrder := models.Order{
			User_Id:    idUser,
			Product_Id: request.Product_Id,
			Qty:        1,
		}
		orderAdded, err := h.OrderRepository.Create_Order(newOrder)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		order, err := h.OrderRepository.GetOrder(orderAdded.Id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{
			Data: convertOrderResponse(order),
		}
		json.NewEncoder(w).Encode(response)
		return

	}
	order.Qty = order.Qty + 1

	orderUpdated, err := h.OrderRepository.UpdateOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	order, err = h.OrderRepository.GetOrder(orderUpdated.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOrderResponse(order),
	}
	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOrderResponse(order)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerOrder) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengabil event dari request body
	var request dto.Update_Request_Order
	json.NewDecoder(r.Body).Decode(&request)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	order, err := h.OrderRepository.GetOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	if request.Event == "add" {
		order.Qty = order.Qty + 1
	} else if request.Event == "less" {
		order.Qty = order.Qty - 1
	}

	if request.Qty != 0 {
		order.Qty = request.Qty
	}

	orderUpdated, err := h.OrderRepository.UpdateOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	order, err = h.OrderRepository.GetOrder(orderUpdated.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOrderResponse(order),
	}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerOrder) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	order, err := h.OrderRepository.GetOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	orderDeleted, err := h.OrderRepository.DeleteOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOrderResponse(orderDeleted),
	}
	json.NewEncoder(w).Encode(res)
}

func convertMultipleOrderResponse(orders []models.Order) []dto.Order_Response {
	var OrderResponse []dto.Order_Response

	for _, order := range orders {
		OrderResponse = append(OrderResponse, dto.Order_Response{
			Id:      order.Id,
			Qty:     order.Qty,
			Product: order.Product,
			User_Id: order.User_Id,
		})
	}

	return OrderResponse
}

func convertOrderResponse(order models.Order) dto.Order_Response {
	return dto.Order_Response{
		Id:      order.Id,
		Qty:     order.Qty,
		Product: order.Product,
	}
}
