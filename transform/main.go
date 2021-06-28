package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/transform/primitive"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func showUploadForm(w http.ResponseWriter, r *http.Request) {
	html := `<html><body>
		<form action="/upload" method="post" enctype="multipart/form-data">
			<input type="file" name="image">
			<button type="submit">Upload Image</button>
		</form>
		</body></html>`

	_, _ = fmt.Fprint(w, html)
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// Get the file extension
	extension := filepath.Ext(header.Filename)[1:]

	output, err := primitive.Transform(file, extension, 10)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch extension {
	case "jpg":
		fallthrough
	case "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case "png":
		w.Header().Set("Content-Type", "image/png")
	default:
		http.Error(w, "Invalid image type", http.StatusBadRequest)
		return
	}

	_, _ = io.Copy(w, output)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", showUploadForm)
	mux.HandleFunc("/upload", uploadImage)

	addr := ":3000"

	fmt.Println(fmt.Sprintf("Server running on port %s", addr))

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
