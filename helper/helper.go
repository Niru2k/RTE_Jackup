package helper

import (
	//User-defined package(s)
	"fmt"
	"online/logs"

	//Third-party package(s)
	"github.com/joho/godotenv"
)

func Config(file string) error {
	log := logs.Log()
	//Load the given file
	if err := godotenv.Load(file); err != nil {
		log.Error.Printf("Error : 'Error at loading %s file'", file)
		fmt.Println("err :", err)
		return err
	}
	return nil
}
