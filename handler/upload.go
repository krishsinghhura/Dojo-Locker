package handler

import (
	"Dojo-Locker/utils"
	"fmt"
	"log"
	"net/http"
	"os"
)

func UploadForHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Println("this method is not allowed")
		return
	}

	file, handler, err := r.FormFile("document")
	if err != nil {
		log.Println("Error retrieving file:", err)
		http.Error(w, "file not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if file == nil || handler == nil {
		log.Println("File or handler is nil")
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}

	//
	tempFile, err := os.CreateTemp("uploads", "uploads-*")
	if err != nil {
		log.Println("Error creating temp file:", err)
		http.Error(w, "Error while creating temp file", http.StatusBadRequest)
	}

	defer tempFile.Close()

	hashString, err := utils.GenerateSha(tempFile.Name())
	if err != nil {
		http.Error(w, "Error generating hash", http.StatusInternalServerError)
		return
	}

	pvtKey, pbcKey, err := utils.GeneratePandPkey(hashString)
	if err != nil {
		log.Println("Error generating keys:", err)
		http.Error(w, "Error generating keys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
        "filename": "%s", 
        "oldHash": "%s",
        "privateKey": "%s",
        "publicKey": "%s"
    }`, handler.Filename, hashString, pvtKey, pbcKey)))
}
