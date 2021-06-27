package primitive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Mode defines the shapes used when transforming the images
type Mode int

// Modes supported by the primitive pkg
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

func primitive(inputFile, outfile string, numberOfShapes int, mode Mode) (string, error) {
	fmt.Println(inputFile, outfile)

	primitiveArgsStr := fmt.Sprintf("-i %s -o %s -n %d -m %d", inputFile, outfile, numberOfShapes, mode)
	primitiveArgs := strings.Fields(primitiveArgsStr)

	cmd := exec.Command("primitive", primitiveArgs...)

	b, err := cmd.CombinedOutput()

	return string(b), err
}

// WithMode returns the Mode to use to transform the image
// Default Mode is ModeTriangle
func WithMode(m Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", m)}
	}
}

func createTempFile(prefix, extension string) (*os.File, error) {
	tempFile, err := ioutil.TempFile("", prefix)

	if err != nil {
		return nil, fmt.Errorf("primitive: failed to create temp input file:: %v", err)
	}

	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile.Name())

	return os.Create(fmt.Sprintf("%s.%s", tempFile.Name(), extension))
}

// Transform takes the provided image and applies primitive transformation to it then returns a reader to the
// resulting image
func Transform(image io.Reader, numberOfShapes int, opts ...func() []string) (io.Reader, error) {
	inputTempFile, err := createTempFile("input_", "png")

	if err != nil {
		return nil, fmt.Errorf("primitive: failed to create temp input file:: %v", err)
	}

	outputTempFile, err := createTempFile("output_", "png")

	if err != nil {
		return nil, fmt.Errorf("primitive: failed to create temp output file:: %v", err)
	}

	// Read image
	_, err = io.Copy(inputTempFile, image)

	if err != nil {
		return nil, fmt.Errorf("primitive: failed to copy temp input file:: %v", err)
	}

	// Run primitive
	stdCombo, err := primitive(inputTempFile.Name(), outputTempFile.Name(), numberOfShapes, ModeCombo)

	if err != nil {
		return nil, fmt.Errorf("primitive: failed to run the primitive command: %v, std combo: %s", err, stdCombo)
	}
	// Read out into a reader, return reader and delete out
	b := bytes.NewBuffer(nil)

	_, err = io.Copy(b, outputTempFile)

	if err != nil {
		return nil, fmt.Errorf("primitive: failed to copy temp output file:: %v", err)
	}

	return b, nil
}
