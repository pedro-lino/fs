package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"fs/pkg/merkle"
	"fs/pkg/types"
)

var rootHash string

// Processes a given directory path and uploads all files within it.
func ProcessPaths(url string, path string, env string) {
	if err := uploadDirectory(url, path, env); err != nil {
		log.Printf("Error uploading directory: %v\n", err)
	}

	if err := os.RemoveAll(path); err != nil {
		log.Printf("Error deleting directory: %v\n", err)
	}
}

// Uploads all files in a directory to the server.
func uploadDirectory(url, dirPath, env string) error {
	var fileHashes []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileHash, err := merkle.HashFile(path)
			if err != nil {
				return fmt.Errorf("error hashing file: %v", err)
			}
			fileHashes = append(fileHashes, fileHash)

			if err := uploadFile(url, path, env); err != nil {
				return err
			}

			if err := os.Remove(path); err != nil {
				return fmt.Errorf("error deleting file: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	mt := merkle.NewMerkleTree()
	for _, hash := range fileHashes {
		mt.AddLeaf(hash)
	}

	rootHash = mt.Root
	if err := saveMerkleRoot([]byte(rootHash)); err != nil {
		return fmt.Errorf("error saving Merkle root: %v", err)
	}

	return nil
}

// Uploads a single file to the server.
func uploadFile(url, path, env string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return fmt.Errorf("error creating form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	client := &http.Client{}

	// Skip certificate verification in dev environment
	if env == "dev" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/upload", url), &body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	return nil
}

func saveMerkleRoot(root []byte) error {
	return os.WriteFile("merkle_root.txt", root, 0644)
}

// Downloads a file from the server and verifies its integrity using the Merkle proof.
func DownloadFile(url, fileName string) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/download/%s", url, fileName), nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error downloading file: bad status: %s\n", resp.Status)
		return
	}

	var response types.DownloadResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response: %v\n", err)
		return
	}

	fileHash := merkle.HashData(response.FileData)
	if !merkle.VerifyProof(fileHash, response.Proof, rootHash) {
		log.Println("File integrity verification failed")
		return
	}

	if err := os.WriteFile(fileName, response.FileData, 0644); err != nil {
		log.Printf("Error saving file: %v\n", err)
		return
	}

	log.Println("File downloaded and verified successfully")
}