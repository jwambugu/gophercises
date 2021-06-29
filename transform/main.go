package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/transform/primitive"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

func createTempFile(prefix, extension string) (*os.File, error) {
	tempFile, err := ioutil.TempFile("./img", prefix)

	if err != nil {
		return nil, fmt.Errorf("main: failed to create temp input file:: %v", err)
	}

	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile.Name())

	return os.Create(fmt.Sprintf("%s.%s", tempFile.Name(), extension))
}

func generateImage(file io.Reader, extension string, numberOfShapes int, mode primitive.Mode) (string, error) {
	output, err := primitive.Transform(file, extension, numberOfShapes, primitive.WithMode(mode))

	if err != nil {
		return "", err
	}

	outputFile, err := createTempFile("", extension)

	_, err = io.Copy(outputFile, output)

	if err != nil {
		return "", err
	}

	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	return outputFile.Name(), nil
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

	a, err := generateImage(file, extension, 33, primitive.ModeCircle)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = file.Seek(0, 0)

	b, err := generateImage(file, extension, 10, primitive.ModeBeziers)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = file.Seek(0, 0)

	html := `<html><body>
		{{ range .}}
			<img src="/{{.}}"> <br/> <br/>
		{{ end }}
	</body></html>`

	tpl := template.Must(template.New("").Parse(html))

	images := []string{a, b}

	if err := tpl.Execute(w, images); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//redirectURL := fmt.Sprintf("/%s", generatedImage)
	//
	//http.Redirect(w, r, redirectURL, http.StatusFound)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", showUploadForm)
	mux.HandleFunc("/upload", uploadImage)

	fs := http.FileServer(http.Dir("./img/"))
	mux.Handle("/img/", http.StripPrefix("/img", fs))

	addr := ":3000"

	fmt.Println(fmt.Sprintf("Server running on port %s", addr))

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
