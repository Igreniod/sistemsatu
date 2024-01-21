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

	fmt.Println("Server is running on port : 8081")
	http.ListenAndServe(":8081", router)
}
