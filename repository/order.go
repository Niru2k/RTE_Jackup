package repository

import (
	//user defined package(s)
	"online/models"

	//Third party package(s)
	"gorm.io/gorm"
)

// Adding a Order into OrderProductInfo table
func CreateOrder(Db *gorm.DB, Order models.OrderProductInfo) error {
	err := Db.Create(&Order).Error
	return err
}

// Delete a Order by Order-id
func DeleteOrderByOrderId(Db *gorm.DB, orderId string) (err error) {
	var Order models.OrderProductInfo
	err = Db.Where("order_id=?", orderId).Delete(&Order).Error
	return
}

// Retrieve orders by user from OrderProductInfo table
func ReadOrdersByUser(Db *gorm.DB, userId string) (Orders []models.OrderProductInfo, err error) {
	err = Db.Where("user_id=?", userId).Find(&Orders).Error
	return
}

// Retrieve orders by admin from OrderProductInfo table
func ReadOrdersByAdmin(Db *gorm.DB) (Orders []models.OrderProductInfo, err error) {
	err = Db.Find(&Orders).Error
	return
}

// Retrieve a Order by Order-id
func ReadOrderByOrderIdUs(Db *gorm.DB, orderId string) (Order models.OrderProductInfo, err error) {
	err = Db.Unscoped().Where("order_id=?", orderId).First(&Order).Error
	return
}

// Retrieve a Order by Order-id
func ReadOrderByOrderId(Db *gorm.DB, orderId string) (Order models.OrderProductInfo, err error) {
	err = Db.Where("order_id=?", orderId).First(&Order).Error
	return
}

// Update a Product by Product-id
func UpdateOrderById(Db *gorm.DB, Order models.OrderProductInfo) (err error) {
	err = Db.Where("order_id=?", Order.OrderId).Save(&Order).Error
	return
}

// Adding a status of the order into Orderstatus table
func CreateOrderStatus(Db *gorm.DB, Order models.OrderStatus) error {
	err := Db.Create(&Order).Error
	return err
}

// Retrieve a max Order-id
func ReadOrderId(Db *gorm.DB) (orderId uint) {
	Db.Model(&models.OrderProductInfo{}).Select("max(order_id)").Scan(&orderId)
	return
}

// Update a status of the order in the Orderstatus table
func UpdateOrderStatus(Db *gorm.DB, Order models.OrderStatus) error {
	err := Db.Where("order_id=?", Order.OrderId).Save(&Order).Error
	return err
}

// Retrieve a OrderStatus by Order-id
func ReadOrderStatusByOrderId(Db *gorm.DB, orderId uint) (Order models.OrderStatus, err error) {
	err = Db.Unscoped().Where("order_id=?", orderId).First(&Order).Error
	return
}

// Retrieve all Order Statuses
func ReadOrderStatus(Db *gorm.DB) (OrderStatus []models.OrderStatus, err error) {
	err = Db.Unscoped().Find(&OrderStatus).Error
	return
}

// Delete a status of the order in the Orderstatus table
func DeleteOrderStatus(Db *gorm.DB, Order models.OrderStatus) error {
	err := Db.Where("order_id=?", Order.OrderId).Delete(&Order).Error
	return err
}
