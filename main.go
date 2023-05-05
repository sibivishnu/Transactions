package main

import (
	"net/http"

	"github.com/sibivishnu/Transactions/routes"
)

func main() {
	router := routes.RegisterRoutes()
	http.ListenAndServe(":8080", router)
}
