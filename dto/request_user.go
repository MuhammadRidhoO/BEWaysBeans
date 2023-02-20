package dto

type Request_Update_Password_User struct {
	Password         string `json:"password" form:"password"`
	Old_Password     string `json:"old_password" form:"old_password"`
	New_Password     string `json:"new_password" form:"new_password"`
	Confirm_Password string `json:"confirm_password" form:"confirm_password"`
}

type Request_Update_User struct {
	Full_Name string `json:"full_name" gorm:"type: varchar(255)"`
	Image     string `json:"image" gorm:"type: varchar(255)"`
}
