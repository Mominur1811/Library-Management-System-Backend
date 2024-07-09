package app

import (
	"fmt"
	"librarymanagement/config"
	"librarymanagement/db"
	"librarymanagement/web"
	"sync"
)

type Application struct {
	wg sync.WaitGroup
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Init() {
	config.LoadConfig()
	db.InitDB()
}

func (app *Application) Run() {
	web.StartServer(&app.wg)
}

func (app *Application) Wait() {
	app.wg.Wait()
}

func (app *Application) CleanUp() {
	db.CloseDB()
}

func HelloWOrld() {
	fmt.Println("HihIHIHIHIH")
}
