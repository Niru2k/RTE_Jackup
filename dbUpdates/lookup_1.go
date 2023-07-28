package dbUpdates

import (
	//user defined package(s)
	"online/driver"
	"online/models"
)

func (Update) Lookup_1() {
	Db := driver.DbConnection()
	Db.AutoMigrate(&models.Roles{})
	Db.AutoMigrate(&models.User{})
	Db.AutoMigrate(&models.Authentication{})
	Db.AutoMigrate(&models.ProductInfo{})
	Db.AutoMigrate(&models.OrderProductInfo{})
	Db.AutoMigrate(&models.OrderStatus{})
	Roles := []models.Roles{
		{RoleId: 1, Role: "admin"},
		{RoleId: 2, Role: "user"},
	}
	Db.Create(&Roles)
}
