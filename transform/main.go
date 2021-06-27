package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/transform/primitive"
	"io"
	"log"
	"os"
)

type PrimitiveMode int

func main() {
	image, err := os.Open("gopher_1.jpg")

	if err != nil {
		log.Fatal(err)
	}

	output, err := primitive.Transform(image, 10)

	if err != nil {
		log.Fatal(err)
	}

	if err := os.Remove("out.png"); err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("out.png")

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(file, output)
	fmt.Println(output)
}
