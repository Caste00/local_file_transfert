package main

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed frontend
var frontendFS embed.FS

func uploadFile(w http.ResponseWriter, req *http.Request) {
	// limite di memoria per il parsing 32 MB
	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join("downloads", header.Filename))
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func listFile(w http.ResponseWriter, req *http.Request) {
	entries, err := os.ReadDir("downloads")
	if err != nil {
		http.Error(w, "Error reading folder", http.StatusInternalServerError)
		return
	}

	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

func download(w http.ResponseWriter, req *http.Request) {
	filename := req.PathValue("filename")
	safe := filepath.Base(filename)
	path := filepath.Join("downloads", safe)

	http.ServeFile(w, req, path)
}

func main() {
	os.MkdirAll("downloads", 0755)

	sub, _ := fs.Sub(frontendFS, "frontend")
	http.Handle("/", http.FileServer(http.FS(sub)))

	http.HandleFunc("POST /upload", uploadFile)
	http.HandleFunc("GET /files", listFile)
	http.HandleFunc("GET /download/{filename}", download)

	http.ListenAndServe(":8092", nil)
}
