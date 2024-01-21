package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//-----run server-----//
	router := mux.NewRouter()
	setRoutes(router)

	fmt.Println("Server is running on port : 8080")
	http.ListenAndServe(":8080", router)

}
