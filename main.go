package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PraveenBSD/content-store/routes"
	"github.com/gorilla/mux"
)

// 	params := mux.Vars(req)
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/upload", routes.UploadContent).Methods("POST")
	router.HandleFunc("/download", routes.DownloadContent).Methods("POST")
	router.HandleFunc("/content-access", routes.ContentAccess).Methods("PUT")
	router.HandleFunc("/info", routes.Info).Methods("GET")
	log.Println("starting server at port 12345")
	log.Fatal(http.ListenAndServe(":12345", router))
	fmt.Println("server running in port 12345")
}
