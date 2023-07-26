package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// Custom Log
type Logs struct {
	Info  *log.Logger
	Error *log.Logger
}

// Login credentials
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup credentials
type SignupReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// User details
type User struct {
	UserId   uint   `json:"-" gorm:"primarykey"`
	Username string `json:"username" binding:"required" gorm:"column:username;type:varchar(100)"`
	Email    string `json:"email" binding:"required" gorm:"column:email;type:varchar(100) unique"`
	Password string `json:"password" binding:"required" gorm:"column:password;type:varchar(100)"`
	Role     string `json:"role" binding:"required" gorm:"-:all"`
	RoleId   uint   `json:"-" gorm:"column:role_id;type:bigint references Roles(role_id)"`
}

// Roles table
type Roles struct {
	RoleId uint   `gorm:"column:role_id;type:bigint primary key"`
	Role   string `gorm:"column:role;type:varchar(50)"`
}

// Token values for each user-id
type Authentication struct {
	UserId uint   `json:"user_id" gorm:"column:user_id;type:bigint primary Key"`
	Token  string `json:"token" gorm:"column:token;type:varchar(200)"`
}

// Credentials for posting a product
type ProductInfoReq struct {
	BrandName    string
	ProductPrice int
	RamCapacity  string
	RamPrice     int
}

// Credentials for posting a product
type OrderProductReq struct {
	BrandName    string `json:"brand_name" binding:"required" `
	ProductPrice int    `json:"product_price" binding:"required"`
	RamCapacity  string `json:"ram_capacity" binding:"required"`
	RamPrice     int    `json:"ram_price" binding:"required"`
	DvdRwDrive   bool   `json:"dvd_rw_drive" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
}

// Details of each product
type ProductInfo struct {
	ProductId    uint   `json:"-" gorm:"primarykey"`
	BrandName    string `json:"brand_name" binding:"required" gorm:"column:brand_name;type:varchar(100)"`
	ProductPrice int    `json:"product_price" binding:"required" gorm:"column:product_price;type:money"`
	RamCapacity  string `json:"ram_capacity" binding:"required" gorm:"column:ram_capacity;type:varchar(100)"`
	RamPrice     int    `json:"ram_price" binding:"required" gorm:"column:ram_price;type:money"`
}

// Details of ordered products
type OrderProductInfo struct {
	OrderId       uint           `json:"-" gorm:"primarykey"`
	ProductId     uint           `json:"-" gorm:"column:product_id;type:bigint references product_infos(product_id)"`
	UserId        uint           `json:"-" gorm:"column:user_id;type:bigint references Users(user_id)"`
	DvdRwDrive    bool           `json:"dvd_rw_drive" binding:"required" gorm:"column:dvd_rw_drive;type:boolean"`
	Name          string         `json:"name" binding:"required" gorm:"column:name;type:varchar(50)"`
	Address       string         `json:"address" binding:"required" gorm:"column:address;type:varchar(200)"`
	PhoneNumber   string         `json:"phone_number" binding:"required" gorm:"column:phone_number;type:varchar(200)"`
	TotalPrice    int            `json:"total_price" binding:"required" gorm:"column:total_price;type:money"`
	PaymentStatus string         `json:"payment_status" binding:"required" gorm:"column:payment_status;type:varchar(50)"`
	CreatedAt     time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"-" gorm:"autoUpdateTime"`
	CancelledAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Tract a order_status
type OrderStatus struct {
	OrderId         uint   `json:"-" gorm:"column:order_id;type:bigint references order_product_infos(order_id)"`
	Name            string `json:"name" gorm:"-"`
	Address         string `json:"address" gorm:"-"`
	PhoneNumber     string `json:"phone_number" gorm:"-"`
	OrderStatus     string `json:"order_status" gorm:"column:order_status;type:varchar(50)"`
	BrandName       string `json:"brand_name" gorm:"-"`
	IncludedProduct string `json:"included_product" gorm:"-"`
	TotalPrice      int    `json:"total_price" gorm:"-"`
}
