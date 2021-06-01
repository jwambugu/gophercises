package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return problems
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a CSV file in the format of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		log.Fatalf("failed to open the CSV file: %s", *csvFilename)
	}

	r := csv.NewReader(file)

	// Parse the CSV file
	records, err := r.ReadAll()

	if err != nil {
		log.Fatalf("failed to parse the CSV: %s", err.Error())
	}

	problems := parseLines(records)

	correctAnswers := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		var answer string
		_, _ = fmt.Scanf("%s\n", &answer)

		if answer == p.answer {
			correctAnswers++
			fmt.Println("Correct!")
		}
	}
	fmt.Printf("You scored %d out of %d\n", correctAnswers, len(problems))
}
