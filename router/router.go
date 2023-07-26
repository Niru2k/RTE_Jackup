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
	handler := handler.Database{Db: Db}
	app.POST("/signup", handler.Signup)
	app.POST("/login", handler.Login)
}

// These handlers are accessible only by admin
func AdminHandlers(Db *gorm.DB, app *echo.Echo) {
	handler := handler.Database{Db: Db}
	middleware := middleware.Database{Db: Db}
	admin := app.Group("/admin", middleware.AuthMiddleware)
	admin.POST("/postProduct", handler.PostProduct)
	admin.PUT("/updateProduct/:product_id", handler.UpdateProductById)
	admin.DELETE("/deleteProduct/:product_id", handler.DeleteProductById)
}

// These handlers are accessible by both admin and user
func CommonHandlers(Db *gorm.DB, app *echo.Echo) {
	handler := handler.Database{Db: Db}
	middleware := middleware.Database{Db: Db}
	common := app.Group("/common", middleware.AuthMiddleware)
	common.GET("/getAllProducts", handler.GetAllProducts, middleware.AuthMiddleware)
}
