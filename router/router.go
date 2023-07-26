package router

import (
	"online/handler"
	"online/logs"
	"online/middleware"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func Router(Db *gorm.DB) {
	log := logs.Log()
	handler := handler.Database{Db: Db}
	middleware := middleware.Database{Db: Db}
	e := echo.New()

	//Common for admin and user
	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)

	//Only for admin
	e.POST("/postproduct", handler.PostProduct, middleware.AuthMiddleware)
	e.POST("/updateproduct/:product_id", handler.PostProduct, middleware.AuthMiddleware)

	//Start a server
	log.Info.Println("Message : 'Server starts in port 8000...' Status : 200")
	if err := e.Start(":8000"); err != nil {
		log.Info.Println("Message : 'Error at start a server...' Status : 500")
		return
	}
}
