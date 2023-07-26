package main

import (
	//user defined package(s)
	"online/driver"
	"online/logs"
	"online/router"

	//Third party package(s)
	"github.com/labstack/echo"
)

func main() {
	log := logs.Log()
	echo := echo.New()

	//Establishing a DB-connection
	Db := driver.DbConnection()

	//Routing all the handlers
	router.LoginHandlers(Db, echo)
	router.AdminHandlers(Db, echo)
	router.CommonHandlers(Db, echo)

	//Start a server
	log.Info.Println("Message : 'Server starts in port 8000...' Status : 200")
	if err := echo.Start(":8000"); err != nil {
		log.Info.Println("Message : 'Error at start a server...' Status : 500")
		return
	}
}
