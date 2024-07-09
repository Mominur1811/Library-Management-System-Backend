package main

import (
	"librarymanagement/app"
)

func main() {

	app := app.NewApplication()
	app.Init()
	app.Run()
	app.Wait()
	app.CleanUp()

}
