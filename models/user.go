package models

type User struct {
	Id               int    `json:"id"`
	Full_Name        string `json:"full_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Old_Password     string `json:"old_password"`
	New_Password     string `json:"new_password"`
	Confirm_Password string `json:"confirm_password"`
	Roles            string `json:"roles"`
	Image            string `json:"image"`
}
type Response_User struct {
	Id               int    `json:"id"`
	Full_Name        string `json:"full_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Old_Password     string `json:"old_password"`
	New_Password     string `json:"new_password"`
	Confirm_Password string `json:"confirm_password"`
	Roles            string `json:"roles"`
	Image            string `json:"image"`
}

func (Response_User) TableName() string {
	return "users"
}
