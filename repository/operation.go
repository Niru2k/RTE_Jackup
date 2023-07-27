package repository

import (
	//user defined package(s)
	"online/models"

	//Third party package(s)
	"gorm.io/gorm"
)

// Adding a specified roles to the roles table
func AddRoles(Db *gorm.DB) {
	Roles := []models.Roles{
		{RoleId: 1, Role: "admin"},
		{RoleId: 2, Role: "user"},
	}
	Db.Create(&Roles)
}

// Table creation
func TableCreation(Db *gorm.DB) {
	Db.AutoMigrate(&models.Roles{})
	Db.AutoMigrate(&models.User{})
	Db.AutoMigrate(&models.Authentication{})
	Db.AutoMigrate(&models.ProductInfo{})
	Db.AutoMigrate(&models.OrderProductInfo{})
	Db.AutoMigrate(&models.OrderStatus{})
	AddRoles(Db)
}
