package main

import (
	"fmt"
	"net/http"
	"login/routes"
	"log"
)

func main(){
	fmt.Println("Hello, World!")
	r := routes.Router()
	fmt.Println("Server started on :5000")
	err := http.ListenAndServe(":5000",r)
	if err != nil {
		log.Fatal(err)
	}

}