package handler

import (
	"Dojo-Locker/utils"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

	tempFile, err := os.CreateTemp("uploads", "uploads-*")
	if err != nil {
		log.Println("Error creating temp file:", err)
		http.Error(w, "Error while creating temp file", http.StatusBadRequest)
	}

	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	reopenedFile, err := os.Open(tempFile.Name())
	if err != nil {
		http.Error(w, "Error opening the file for hashing", http.StatusInternalServerError)
		return
	}
	defer reopenedFile.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, reopenedFile); err != nil {
		http.Error(w, "Error generating hash", http.StatusInternalServerError)
		return
	}

	hashString := hex.EncodeToString(hash.Sum(nil))

	fmt.Printf("File: %s | hash: %s\n", handler.Filename, hashString)

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
