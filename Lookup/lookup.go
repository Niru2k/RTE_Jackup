package Lookup

import (
	//user defined package(s)
	"online/dbUpdates"
	"online/models"

	//Inbuild package(s)
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	//Third party package(s)
	"gorm.io/gorm"
)

func UpdateDatabase(Db *gorm.DB) {
	var update []models.Updates
	var check bool
	files, err := os.ReadDir("./dbUpdates")
	if err != nil {
		log.Println("Error :", err)
		return
	}
	Db.AutoMigrate(&models.Updates{})
	Db.Find(&update)
	for i := 0; i < len(files); i++ {
		for _, value := range update {
			file := fmt.Sprintf("%s.go", value.FileName)
			if files[i].Name() == file {
				check = true
				break
			}
		}
		if !check && files[i].Name() != "update.go" {
			var update1 models.Updates
			extension := path.Ext(files[i].Name())
			index := strings.LastIndex(files[i].Name(), extension)
			update1.FileName = files[i].Name()[:index]
			Db.Create(&update1)
			log.Println(update1.FileName, "is updated...")
			update := dbUpdates.Update{}
			update.Invoke(update1.FileName)
		}
		check = false
	}

}
