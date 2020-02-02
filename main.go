package main

import (
	"github.com/Deewai/finleap/app"
	"os"
)

func main() {
	a := app.App{}
	// Make sure environment variables are set
	a.Initialize(os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE"))

	a.Run(":3000")
}
