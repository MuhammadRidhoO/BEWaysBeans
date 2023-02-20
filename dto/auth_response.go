package dto

type Response_Register struct {
	Id        int    `json:"id"`
	Full_Name string `json:"full_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Roles     string `json:"roles"`
}
type Response_Login struct {
	Id        int    `json:"id"`
	Full_Name string `json:"full_name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	Roles     string `json:"roles"`
}

type Response_CheckOut struct {
	Id        int    `json:"id"`
	Full_Name string `json:"full_name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	Roles     string `json:"roles"`
}
