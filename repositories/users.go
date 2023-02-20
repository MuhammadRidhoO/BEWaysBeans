package repositories

import (
	"BEWaysBeans/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUsers() ([]models.User, error)
	GetUser(ID int) (models.User, error)
	GetUserLogin(ID int) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	UpdatePasswordUser(update_password_User models.User) (models.User, error)
	DeleteUser(user models.User) (models.User, error)

	// GetPassword(username string) (string, error)
	// UpdatePassword(username, password string) error
}

func RepositoryUser(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindUsers() ([]models.User, error) {
	var user []models.User
	err := r.db.Find(&user).Error

	return user, err
}

func (r *repository) GetUser(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

func (r *repository) GetUserLogin(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

func (r *repository) UpdateUser(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error

	return user, err
}
func (r *repository) UpdatePasswordUser(update_password_User models.User) (models.User, error) {
	err := r.db.Save(&update_password_User).Error

	return update_password_User, err
}
func (r *repository) DeleteUser(user models.User) (models.User, error) {
	err := r.db.Delete(&user).Error

	return user, err
}
