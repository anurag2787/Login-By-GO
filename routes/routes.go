package routes

import (
	"net/http"
	"github.com/gorilla/mux"

	control "login/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/login", control.LoginUser).Methods("POST")
	router.HandleFunc("/api/register", control.RegisterUser).Methods("POST")
	router.HandleFunc("/api",func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type","application/json")
		w.Write([]byte(`{"Server Started":true}`))
	}).Methods("GET")

	return router	
}