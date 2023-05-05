package routes

import (
	"github.com/gorilla/mux"
	"github.com/sibivishnu/Transactions/controllers"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/transactions", controllers.PostTransaction).Methods("POST")
	router.HandleFunc("/transactions", controllers.DeleteTransactions).Methods("DELETE")
	router.HandleFunc("/statistics", controllers.GetStatistics).Methods("GET")
	router.HandleFunc("/location", controllers.SetLocation).Methods("POST")
	router.HandleFunc("/location/reset", controllers.ResetLocation).Methods("POST")

	return router
}
