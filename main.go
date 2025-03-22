package main

import (
	"Dojo-Locker/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/upload", handler.UploadForHash)
	fmt.Println("Server running on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
