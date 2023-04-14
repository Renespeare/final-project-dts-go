package main

import (
	"final-project-dts-go/database"
	"final-project-dts-go/router"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main()  {
	// var PORT = ":8080"
	database.StartDB()

	fmt.Println("Successfully connected to database")

	// routers.StartServer().Run(PORT)
	routers.StartServer().Run(":" + os.Getenv("PORT"))
}