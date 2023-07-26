package repository

import (
	//user defined package(s)
	"online/models"

	//Third party package(s)
	"gorm.io/gorm"
)

// Adding a token into authorizations table
func AddToken(Db *gorm.DB, auth models.Authentication) error {
	err := Db.Create(&auth).Error
	return err
}

// Delete a token by user-id
func DeleteToken(Db *gorm.DB, userId string) (err error) {
	var token models.Authentication
	err = Db.Where("user_id=?", userId).Delete(&token).Error
	return
}
