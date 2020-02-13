package main

import (
	"app/config"
	app "app/src"
	"log"
)

func main() {
	_config := config.GetConfig()

	_app := &app.App{}
	_app.SetConfig(_config)
	if _app.SetDB() {
		_app.LoadRoutes()
		_app.InitApp()
	}
	log.Fatal("Error connect with DB")
}
