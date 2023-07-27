package router

import (
	//user defined packages
	"online/handler"
	"online/middleware"

	//Third party packages
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// Signup and Login Handlers
func LoginHandlers(Db *gorm.DB, app *echo.Echo) {
	handler := handler.Database{Connection: Db}
	app.POST("/signup", handler.Signup)
	app.POST("/login", handler.Login)
}

// These handlers are accessible only by admin
func AdminHandlers(Db *gorm.DB, app *echo.Echo) {
	handler := handler.Database{Connection: Db}
	middleware := middleware.Database{Connection: Db}
	admin := app.Group("/admin", middleware.AuthMiddleware)
	admin.POST("/postProduct", handler.PostProduct)
	admin.PUT("/updateProduct/:product_id", handler.UpdateProductById)
	admin.DELETE("/deleteProduct/:product_id", handler.DeleteProductById)
	admin.PUT("/updateStatus/:order_id", handler.UpdateOrderStatusById)
}

// These handlers are accessible only by user
func UserHandlers(Db *gorm.DB, app *echo.Echo) {
	handler := handler.Database{Connection: Db}
	middleware := middleware.Database{Connection: Db}
	user := app.Group("/user", middleware.AuthMiddleware)
	user.POST("/postOrder", handler.AddOrder)
	user.DELETE("/cancelOrder/:order_id", handler.CancelOrderById)
	user.POST("/payment/:order_id", handler.Payment)
}

// These handlers are accessible by both admin and user
func CommonHandlers(Db *gorm.DB, app *echo.Echo) {
	handler := handler.Database{Connection: Db}
	middleware := middleware.Database{Connection: Db}
	common := app.Group("/common", middleware.AuthMiddleware)
	common.GET("/getAllProducts", handler.GetAllProducts)
	common.GET("/getOrders", handler.GetOrders)
	common.GET("/getOrderStatus/:order_id", handler.GetOrderStatusById)
}
