// main.go

package main

import (
	"basictrade/routes"
	database "basictrade/utils"

)

var (
	PORT = ":5050"
)

func main() {
	database.StartDB()
	r := routes.StartApp()
	r.Run(PORT)
}
