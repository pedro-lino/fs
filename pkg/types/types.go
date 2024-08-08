package types

type Url struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

type FileUploadRequest struct {
    FileName string `json:"fileName"`
    FileData []byte `json:"fileData"`
}

type UploadResponse struct {
	LeafHashes []string `json:"leaf_hashes"`
}

type DownloadResponse struct {
	FileData []byte   `json:"file_data"`
	Proof    []string `json:"proof"`
}