package repository

import (
	//user defined package(s)
	"online/models"

	//Third party package(s)
	"gorm.io/gorm"
)

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

// Delete a product by product-id
func DeleteProductByProductId(Db *gorm.DB, ProductId string) (err error) {
	var Product models.ProductInfo
	err = Db.Where("Product_id=?", ProductId).Delete(&Product).Error
	return
}

// Retrieve all products from products table
func ReadAllProducts(Db *gorm.DB) (Products []models.ProductInfo, err error) {
	err = Db.Find(&Products).Error
	return
}

// Retrieve product by products's specifications
func ReadProductIdByProductData(Db *gorm.DB, Product models.OrderProductInfo) (product models.ProductInfo, err error) {
	err = Db.Select("product_id").Where("brand_name=? AND product_price=? AND ram_capacity=? AND ram_price=?", Product.BrandName, Product.ProductPrice, Product.RamCapacity, Product.RamPrice).First(&product).Error
	return
}

// Retrieve products's specifications by product_id
func ReadProductDataByProductId(Db *gorm.DB, ProductId uint) (product models.ProductInfo, err error) {
	err = Db.Where("product_id=?", ProductId).First(&product).Error
	return
}
