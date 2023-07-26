package repository

import (
	//user defined package(s)
	"online/models"

	//Third party package(s)
	"gorm.io/gorm"
)

// Retrieve the User's role by role-id
func ReadRoleIdByRole(Db *gorm.DB, data models.User) (models.Roles, error) {
	role := models.Roles{}
	err := Db.Select("role_id").Where("role=?", data.Role).First(&role).Error
	return role, err
}

// Adding a user details into users table
func CreateUser(Db *gorm.DB, data models.User) (err error) {
	err = Db.Create(&data).Error
	return
}

// Retrieve the User details by Email
func ReadUserByEmail(Db *gorm.DB, data models.User) (models.User, error) {
	err := Db.Where("email = ?", data.Email).First(&data).Error
	return data, err
}

// Retrieve a token by user-id
func ReadTokenByUserId(Db *gorm.DB, user models.User) (auth models.Authentication, err error) {
	err = Db.Where("user_id=?", user.UserId).First(&auth).Error
	return auth, err
}
