package handlers

import (
	"encoding/json"
	"fs/pkg/merkle"
	"fs/pkg/types"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

var MerkleTree *merkle.MerkleTree

func init() {
	MerkleTree = merkle.NewMerkleTree()
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		log.Println("Error parsing form data:", err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file from form data", http.StatusBadRequest)
		log.Println("Error retrieving file:", err)
		return
	}
	defer file.Close()

	filePath := filepath.Join("./uploaded_files", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating destination file", http.StatusInternalServerError)
		log.Println("Error creating file:", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		log.Println("Error saving file:", err)
		return
	}

	// Hash the uploaded file and add it to the Merkle tree
	fileHash, err := merkle.HashFile(filePath)
	if err != nil {
		http.Error(w, "Error hashing file", http.StatusInternalServerError)
		log.Println("Error hashing file:", err)
		return
	}
	MerkleTree.AddLeaf(fileHash)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	filePath := filepath.Join("./uploaded_files", fileName)

	info, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Println("Error reading file:", err)
		return
	}

	if info.IsDir() {
		http.Error(w, "Requested file is a directory", http.StatusBadRequest)
		log.Println("Error: Requested file is a directory")
		return
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		log.Println("Error reading file:", err)
		return
	}

	fileHash:= merkle.HashData(fileData)

	proof, err := MerkleTree.GenerateProof(fileHash)
	if err != nil {
		log.Println("Error generating proof, sending empty proof:", err)
		proof = []string{}
	}

	response := types.DownloadResponse{
		FileData: fileData,
		Proof:    proof,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}