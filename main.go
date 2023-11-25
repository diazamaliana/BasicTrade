// main.go

package main

import (
	"basictrade/routes"
	database "basictrade/utils"
	"os"
)


func main() {
	// Start the database connection
	database.StartDB()

	// Get the port from the environment variable or use a default value
	port := os.Getenv("PORT")
	if port == "" {
		port = "5050" // Default port if PORT environment variable is not set
	}

	// Start the application on the specified port
	r := routes.StartApp()
	r.Run(":" + port)
}
