package dto

type Response_User struct {
	Id               int    `json:"id"`
	Full_Name        string `json:"full_name" gorm:"type: varchar(255)"`
	Email            string `json:"email" gorm:"type: varchar(255)"`
	Roles            string `json:"roles" form:"roles"`
	Password         string `json:"password" form:"password"`
	// Old_Password         string `json:"old_password" form:"old_password"`
	// New_Password     string `json:"new_password" form:"new_password"`
	Confirm_Password string `json:"confirm_password" form:"confirm_password"`
	Image            string `json:"image" form:"image"`
}

type Response_Update_User struct {
	Id        int    `json:"id"`
	Full_Name string `json:"full_name" gorm:"type: varchar(255)"`
	Email     string `json:"email" gorm:"type: varchar(255)"`
	Image     string `json:"image" form:"image"`
	Roles     string `json:"roles" form:"roles"`
}

type Response_Update_Password_User struct {
	Id               int    `json:"id"`
	Full_Name        string `json:"full_name" gorm:"type: varchar(255)"`
	Email            string `json:"email" gorm:"type: varchar(255)"`
	Roles            string `json:"roles" form:"roles"`
	Password         string `json:"password" form:"password"`
	New_Password     string `json:"new_password" form:"new_password"`
	Confirm_Password string `json:"confirm_password" form:"confirm_password"`
}
