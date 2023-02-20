package dto

type Request_Register struct {
	Full_Name string `json:"full_name" gorm:"type: varchar(255)" validate:"required"`
	Email     string `json:"email" gorm:"type: varchar(255)" validate:"required"`
	Password  string `json:"password" gorm:"type: varchar(255)" validate:"required"`
	Roles   string `json:"roles" gorm:"type: varchar(255)"`
}

type Request_Login struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}
