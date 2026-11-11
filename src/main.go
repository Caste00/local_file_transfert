package main

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"fmt"
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

func deleteFile(w http.ResponseWriter, req *http.Request) {
	filename := req.PathValue("filename")
	safe := filepath.Base(filename)
	path := filepath.Join("downloads", safe)
	
	err := os.Remove(path)
	if err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func download(w http.ResponseWriter, req *http.Request) {
	filename := req.PathValue("filename")
	safe := filepath.Base(filename)
	path := filepath.Join("downloads", safe)

	http.ServeFile(w, req, path)
}

func printAddresses(port string) {
    interfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("Error retrieving network interfaces:", err)
        return
    }

    for _, iface := range interfaces {
        // Ignora interfacce non attive e loopback
        if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
            continue
        }

        addresses, err := iface.Addrs()
        if err != nil {
            continue
        }

        for _, addr := range addresses {
            ipNet, ok := addr.(*net.IPNet)
            if !ok || ipNet.IP.To4() == nil {
                continue
            }

            fmt.Printf("http://%s:%s\n", ipNet.IP.String(), port)
        }
    }
}

func main() {
	os.MkdirAll("downloads", 0755)

	sub, _ := fs.Sub(frontendFS, "frontend")
	http.Handle("/", http.FileServer(http.FS(sub)))

	http.HandleFunc("POST /upload", uploadFile)
	http.HandleFunc("GET /files", listFile)
	http.HandleFunc("DELETE /files/{filename}", deleteFile)
	http.HandleFunc("GET /download/{filename}", download)

	fmt.Println("Running on:")
	printAddresses("8090")

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
