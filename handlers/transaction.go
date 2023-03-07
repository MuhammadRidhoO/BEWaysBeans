package handlers

import (
	"BEWaysBeans/dto"
	"BEWaysBeans/models"
	"BEWaysBeans/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

// Create `path_file` Global variable here ...

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transaction, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: transaction}
	json.NewEncoder(w).Encode(response)
	// fmt.Println(transaction)
}

func (h *handlerTransaction) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Id := mux.Vars(r)["id"]

	var Transaction models.Transaction
	Transaction, err := h.TransactionRepository.GetTransaction(Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create Embed Path File on Image property here ...
	// Transaction.Image_Payment = os.Getenv("PATH_FILE") + Transaction.Image_Payment

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: Transaction}
	json.NewEncoder(w).Encode(response)
}
func (h *handlerTransaction) FindTransactionsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	idUser := int(userInfo["id"].(float64))

	transactions, err := h.TransactionRepository.FindTransactionsByUser(idUser)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: transactions,
	}
	json.NewEncoder(w).Encode(res)
}
func (h *handlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request dto.CreateTransactionRequest
	json.NewDecoder(r.Body).Decode(&request)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	request.User_Id = int(userInfo["id"].(float64))

	// memvalidasi inputan dari request body
	validation := validator.New()
	errValidation := validation.Struct(request)
	if errValidation != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "Failed Validation Transaction"}
		json.NewEncoder(w).Encode(response)
		return
	}
	// var TransIdIsMatch = false
	// var TransactionId int
	// for !TransIdIsMatch {
	// 	TransactionId = int(time.Now().Unix()) // 12948129048123
	// 	transactionData, _ := h.TransactionRepository.GetTransaction(TransactionId)
	// 	if transactionData.Id == 0 {
	// 		TransIdIsMatch = true
	// 	}
	// }
	newTransaction := models.Transaction{
		Id:              fmt.Sprintf("TRX-%d-%d", request.User_Id, timeIn("Asia/Jakarta").UnixNano()),
		Total:          request.Total,
		OrderDate:      timeIn("Asia/Jakarta"),
		Status_Payment: "pending",
		User_Id:        request.User_Id,
	}

	for _, order := range request.Products {
		newTransaction.Order = append(newTransaction.Order, models.Order_Response_For_Transaction{
			Id:         order.Id,
			Product_Id: order.Product_Id,
			Qty:        order.Qty,
		})
	}

	transactionAdded, err := h.TransactionRepository.CreateTransaction(newTransaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: "Failed Create Transaction",
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	transaction, err := h.TransactionRepository.GetTransaction(transactionAdded.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: "Transaction Not Found",
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	// 2. Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.Id,
			GrossAmt: int64(transaction.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transaction.User.Full_Name,
			Email: transaction.User.Email,
		},
	}

	snapResp, err := s.CreateTransaction(req)
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: snapResp}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var request dto.UpdateTransactionRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	_, err = h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	// mengupdate status transaksi
	transaction, err := h.TransactionRepository.UpdateTransaction(request.Status_Payment, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	// mengambil data transaksi yang sudah diupdate
	transaction, err = h.TransactionRepository.GetTransaction(transaction.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: transaction}
	json.NewEncoder(w).Encode(response)
}

func SendMail(status_payment string, transaction models.Transaction) {

	//  if status != transaction.Status && status == "success" {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "Muhammad Ridho <muhammadridho081997@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
	var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")
	//  }
	var products []map[string]interface{}

	for _, order := range transaction.Order {
		products = append(products, map[string]interface{}{
			"Name_Product": order.Product.Name_Product,
			"Price":        order.Product.Price,
			"Qty":          order.Qty,
			"Total":        order.Qty * order.Product.Price,
		})
	}

	data := map[string]interface{}{
		"Transaction_Id": transaction.Id,
		"Status_Payment": status_payment,
		"Full_Name":      transaction.User.Full_Name,
		"Total":          transaction.Total,
		"Products":       products,
	}
	t, err := template.ParseFiles("view/notification_email.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bodyMail := new(bytes.Buffer)

	fmt.Println(bodyMail)

	// mengeksekusi template, dan memparse "data" ke template
	t.Execute(bodyMail, data)

	name_product := products
	price_product := strconv.Itoa(transaction.Total)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", transaction.User.Email)
	mailer.SetHeader("Subject", "Transaction Status Payment")
	mailer.SetBody("text/html", fmt.Sprintf(`<!doctype html>
	<html>
	  <head>
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
		<title>Simple Transactional Email</title>
		<style>
		  img {
			border: none;
			-ms-interpolation-mode: bicubic;
			max-width: 100vw; 
		  }
	
		  body {
			background-color: #f6f6f6;
			font-family: sans-serif;
			-webkit-font-smoothing: antialiased;
			font-size: 14px;
			line-height: 1.4;
			margin: 0;
			padding: 0;
			-ms-text-size-adjust: 100vw;
			-webkit-text-size-adjust: 100vw; 
		  }
	
		  table {
			border-collapse: separate;
			mso-table-lspace: 0pt;
			mso-table-rspace: 0pt;
			width: 100vw; }
			table td {
			  font-family: sans-serif;
			  font-size: 14px;
			  vertical-align: top; 
		  }
	
		  .body {
			background-color: #f6f6f6;
			width: 100vw; 
		  }
	
		  .container {
			display: block;
			margin: 0 auto !important;
			/* makes it centered */
			max-width: 580px;
			padding: 10px;
			width: 580px; 
		  }
	
		  .content {
			box-sizing: border-box;
			display: block;
			margin: 0 auto;
			max-width: 580px;
			padding: 10px; 
		  }

		  .main {
			background: #ffffff;
			border-radius: 3px;
			width: 100vw; 
		  }
	
		  .wrapper {
			box-sizing: border-box;
			padding: 20px; 
		  }
	
		  .content-block {
			padding-bottom: 10px;
			padding-top: 10px;
		  }
	
		  .footer {
			clear: both;
			margin-top: 10px;
			text-align: center;
			width: 100vw; 
		  }
			.footer td,
			.footer p,
			.footer span,
			.footer a {
			  color: #999999;
			  font-size: 12px;
			  text-align: center; 
		  }
	
		  h1,
		  h2,
		  h3,
		  h4 {
			color: #000000;
			font-family: sans-serif;
			font-weight: 400;
			line-height: 1.4;
			margin: 0;
			margin-bottom: 30px; 
		  }
	
		  h1 {
			font-size: 35px;
			font-weight: 300;
			text-align: center;
			text-transform: capitalize; 
		  }
	
		  p,
		  ul,
		  ol {
			font-family: sans-serif;
			font-size: 14px;
			font-weight: normal;
			margin: 0;
			margin-bottom: 15px; 
		  }
			p li,
			ul li,
			ol li {
			  list-style-position: inside;
			  margin-left: 5px; 
		  }
	
		  a {
			color: #3498db;
			text-decoration: underline; 
		  }
	
		  .btn {
			box-sizing: border-box;
			width: 100vw; }
			.btn > tbody > tr > td {
			  padding-bottom: 15px; }
			.btn table {
			  width: auto; 
		  }
			.btn table td {
			  background-color: #ffffff;
			  border-radius: 5px;
			  text-align: center; 
		  }
			.btn a {
			  background-color: #ffffff;
			  border: solid 1px #3498db;
			  border-radius: 5px;
			  box-sizing: border-box;
			  color: #3498db;
			  cursor: pointer;
			  display: inline-block;
			  font-size: 14px;
			  font-weight: bold;
			  margin: 0;
			  padding: 12px 25px;
			  text-decoration: none;
			  text-transform: capitalize; 
		  }
	
		  .btn-primary table td {
			background-color: #3498db; 
		  }
	
		  .btn-primary a {
			background-color: #3498db;
			border-color: #3498db;
			color: #ffffff; 
		  }

		  .last {
			margin-bottom: 0; 
		  }
	
		  .first {
			margin-top: 0; 
		  }
	
		  .align-center {
			text-align: center; 
		  }
	
		  .align-right {
			text-align: right; 
		  }
	
		  .align-left {
			text-align: left; 
		  }
	
		  .clear {
			clear: both; 
		  }
	
		  .mt0 {
			margin-top: 0; 
		  }
	
		  .mb0 {
			margin-bottom: 0; 
		  }
	
		  .preheader {
			color: transparent;
			display: none;
			height: 0;
			max-height: 0;
			max-width: 0;
			opacity: 0;
			overflow: hidden;
			mso-hide: all;
			visibility: hidden;
			width: 0; 
		  }
	
		  .powered-by a {
			text-decoration: none; 
		  }
	
		  hr {
			border: 0;
			border-bottom: 1px solid #f6f6f6;
			margin: 20px 0; 
		  }
	
		  @media only screen and (max-width: 620px) {
			table.body h1 {
			  font-size: 28px !important;
			  margin-bottom: 10px !important; 
			}
			table.body p,
			table.body ul,
			table.body ol,
			table.body td,
			table.body span,
			table.body a {
			  font-size: 16px !important; 
			}
			table.body .wrapper,
			table.body .article {
			  padding: 10px !important; 
			}
			table.body .content {
			  padding: 0 !important; 
			}
			table.body .container {
			  padding: 0 !important;
			  width: 100vw !important; 
			}
			table.body .main {
			  border-left-width: 0 !important;
			  border-radius: 0 !important;
			  border-right-width: 0 !important; 
			}
			table.body .btn table {
			  width: 100vw !important; 
			}
			table.body .btn a {
			  width: 100vw !important; 
			}
			table.body .img-responsive {
			  height: auto !important;
			  max-width: 100vw !important;
			  width: auto !important; 
			}
		  }
	
		  /* -------------------------------------
			  PRESERVE THESE STYLES IN THE HEAD
		  ------------------------------------- */
		  @media all {
			.ExternalClass {
			  width: 100vw; 
			}
			.ExternalClass,
			.ExternalClass p,
			.ExternalClass span,
			.ExternalClass font,
			.ExternalClass td,
			.ExternalClass div {
			  line-height: 100vw; 
			}
			.apple-link a {
			  color: inherit !important;
			  font-family: inherit !important;
			  font-size: inherit !important;
			  font-weight: inherit !important;
			  line-height: inherit !important;
			  text-decoration: none !important; 
			}
			#MessageViewBody a {
			  color: inherit;
			  text-decoration: none;
			  font-size: inherit;
			  font-family: inherit;
			  font-weight: inherit;
			  line-height: inherit;
			}
			.btn-primary table td:hover {
			  background-color: #34495e !important; 
			}
			.btn-primary a:hover {
			  background-color: #34495e !important;
			  border-color: #34495e !important; 
			} 
		  }
	
		</style>
	  </head>
	  <body>
		<span class="preheader">This is preheader text. Some clients will show this text as a preview.</span>
		<table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body">
		  <tr>
			<td>&nbsp;</td>
			<td class="container">
			  <div class="content">
	
				<!-- START CENTERED WHITE CONTAINER -->
				<table role="presentation" class="main">
	
				  <!-- START MAIN CONTENT AREA -->
				  <tr>
					<td class="wrapper">
					  <table role="presentation" border="0" cellpadding="0" cellspacing="0">
						<tr>
						  <td>
							<p>Hi there,</p>
							<p>Sometimes you just want to send a simple HTML email with a simple design and clear call to action. This is it.</p>
							<table role="presentation" border="0" cellpadding="0" cellspacing="0" class="btn btn-primary">
							  <tbody>
								<tr>
								  <td align="left">
									<table role="presentation" border="0" cellpadding="0" cellspacing="0">
									  <tbody>
										<tr>
										  <td> Product Name : %s </td>
										  <td> Price : %s </td>
										  <td> Status : %s </td>
										</tr>
									  </tbody>
									</table>
								  </td>
								</tr>
							  </tbody>
							</table>
							<p>This is a really simple email template. Its sole purpose is to get the recipient to click the button with no distractions.</p>
							<p>Good luck! Hope it works.</p>
						  </td>
						</tr>
					  </table>
					</td>
				  </tr>
	
				<!-- END MAIN CONTENT AREA -->
				</table>
				<!-- END CENTERED WHITE CONTAINER -->
	
			  </div>
			</td>
			<td>&nbsp;</td>
		  </tr>
		</table>
	  </body>
	</html>`, name_product, price_product, status_payment))

	dialer := gomail.NewDialer(CONFIG_SMTP_HOST, CONFIG_SMTP_PORT, CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent! to " + transaction.User.Email)
}

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	// 1. Initialize empty map
	var notificationPayload map[string]interface{}

	// 2. Parse JSON request body and use it to set json to payload
	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		// do something on error when decode
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	// 3. Get order-id from payload
	orderId, _ := notificationPayload["order_id"].(string)
	// if !exists {
	// 	// do something when key `order_id` not found
	// 	return
	// }

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)

	transaction, err := h.TransactionRepository.GetTransactionString(orderId)

	if err != nil {
		fmt.Println("Transaction not found")
		return
	}

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			SendMail("success", transaction)
			h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
		} else if fraudStatus == "accept" {
			SendMail("success", transaction)
			h.TransactionRepository.UpdateTransaction("success", transaction.Id)
		}
	} else if transactionStatus == "settlement" {
		SendMail("success", transaction)
		h.TransactionRepository.UpdateTransaction("success", transaction.Id)
		fmt.Println("test 3")
	} else if transactionStatus == "deny" {
		SendMail("success", transaction)
		h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		SendMail("success", transaction)
		h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
	} else if transactionStatus == "pending" {
		SendMail("success", transaction)
		h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
	}

}

func timeIn(name string) time.Time {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}
