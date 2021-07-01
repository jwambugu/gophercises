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
	"strconv"
	"time"
)

type (
	generatedImage struct {
		Name           string
		Mode           primitive.Mode
		NumberOfShapes int
	}

	generateImagesOptions struct {
		NumberOfShapes int
		Mode           primitive.Mode
	}
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

	fileOnDisk, err := createTempFile("", extension)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func(fileOnDisk *os.File) {
		_ = fileOnDisk.Close()
	}(fileOnDisk)

	_, err = io.Copy(fileOnDisk, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/modify/%s", filepath.Base(fileOnDisk.Name()))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func imageHeadersToSet(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".jpeg":
		fallthrough
	case ".jpg":
		return "image/jpg"
	default:
		return ""
	}
}

func generateImages(rs io.ReadSeeker, extension string, opts ...generateImagesOptions) ([]generatedImage, error) {
	var images []generatedImage

	for _, opt := range opts {
		_, _ = rs.Seek(0, 0)

		image, err := generateImage(rs, extension, opt.NumberOfShapes, opt.Mode)

		if err != nil {
			return nil, err
		}

		fmt.Println(fmt.Sprintf("generated image: %s", image))

		images = append(images, generatedImage{
			Name:           filepath.Base(image),
			Mode:           opt.Mode,
			NumberOfShapes: opt.NumberOfShapes,
		})
	}

	return images, nil
}

func renderModeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, extension string) {
	modes := []primitive.Mode{
		primitive.ModeCircle,
		primitive.ModeBeziers,
		primitive.ModePolygon,
		primitive.ModeCombo,
	}

	var opts []generateImagesOptions

	for _, mode := range modes {
		opts = append(opts, generateImagesOptions{
			NumberOfShapes: 10,
			Mode:           mode,
		})
	}

	images, err := generateImages(rs, extension, opts...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{ range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}">
					<img style="width: 20%" src="/img/{{.Name}}">
				</a>
			{{ end }}
		</body></html>`

	tpl := template.Must(template.New("").Parse(html))

	if err := tpl.Execute(w, images); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderNumberOfShapesChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, extension string,
	mode primitive.Mode) {

	opts := []generateImagesOptions{
		{NumberOfShapes: 10, Mode: mode},
		{NumberOfShapes: 20, Mode: mode},
		{NumberOfShapes: 30, Mode: mode},
		{NumberOfShapes: 40, Mode: mode},
	}

	images, err := generateImages(rs, extension, opts...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{ range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}&n={{.NumberOfShapes}}">
					<img style="width: 20%" src="/img/{{.Name}}">
				</a>
			{{ end }}
		</body></html>`

	tpl := template.Must(template.New("").Parse(html))

	if err := tpl.Execute(w, images); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fileToModify(w http.ResponseWriter, r *http.Request) {
	imagePath := fmt.Sprintf("./img/%s", filepath.Base(r.URL.Path))

	file, err := os.Open(imagePath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	imageExtension := filepath.Ext(imagePath)[1:]
	//imageHeaders := imageHeadersToSet(imageExtension)

	selectedMode := r.FormValue("mode")

	if selectedMode == "" {
		//	render mode choices
		renderModeChoices(w, r, file, imageExtension)
		return
	}

	mode, err := strconv.Atoi(selectedMode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	selectedNumberOfShapes := r.FormValue("n")

	if selectedNumberOfShapes == "" {
		// render number of shapes choices
		renderNumberOfShapesChoices(w, r, file, imageExtension, primitive.Mode(mode))
		return
	}

	numberOfShapes, err := strconv.Atoi(selectedNumberOfShapes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = numberOfShapes
	//w.Header().Set("Content-Type", imageHeaders)

	_, err = io.Copy(w, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/img/%s", file.Name())

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func main() {

	go func() {
		ticker := time.NewTicker(2 * time.Minute)

		for {
			<-ticker.C
			// TODO: Check for images to delete
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/", showUploadForm)
	mux.HandleFunc("/upload", uploadImage)
	mux.HandleFunc("/modify/", fileToModify)

	fs := http.FileServer(http.Dir("./img/"))
	mux.Handle("/img/", http.StripPrefix("/img", fs))

	addr := ":3000"

	fmt.Println(fmt.Sprintf("Server running on port %s", addr))

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// a, err := generateImage(file, extension, 33, primitive.ModeCircle)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	_, _ = file.Seek(0, 0)
//
//	b, err := generateImage(file, extension, 10, primitive.ModeBeziers)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	_, _ = file.Seek(0, 0)
//
//	html := `<html><body>
//		{{ range .}}
//			<img src="/{{.}}"> <br/> <br/>
//		{{ end }}
//	</body></html>`
//
//	tpl := template.Must(template.New("").Parse(html))
//
//	images := []string{a, b}
//
//	if err := tpl.Execute(w, images); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
