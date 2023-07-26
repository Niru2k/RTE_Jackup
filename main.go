package main

import (
	"online/driver"
	"online/router"
)

func main() {
	//Establishing a DB-connection
	Db := driver.DbConnection()

	//Routing all the handlers
	router.Router(Db)
}
