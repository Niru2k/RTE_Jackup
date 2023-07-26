package repository

import (
	"online/models"

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

// Adding a product into products table
func CreateProduct(Db *gorm.DB, Product models.ProductInfo) error {
	err := Db.Create(&Product).Error
	return err
}

// Retrieve a product by product-id
func ReadProductByProductId(Db *gorm.DB, productId string) (product models.ProductInfo, err error) {
	err = Db.Where("product_id=?", productId).First(&product).Error
	return
}

// Update a Product by Product-id
func UpdateProductByProductId(Db *gorm.DB, ProductId string, Product models.ProductInfo) (err error) {
	err = Db.Where("product_id=?", ProductId).Save(&Product).Error
	return
}
