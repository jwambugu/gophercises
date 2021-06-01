package main

import (
	"fmt"
	"strings"
)

// There is a sequence of words in CamelCase as a string of letters, , having the following properties:
// It is a concatenation of one or more words consisting of English letters.
// All letters in the first word are lowercase.
// For each of the subsequent words, the first letter is uppercase and rest of the letters are lowercase.
// Given , determine the number of words in .

// Example
// There are  words in the string: 'one', 'Two', 'Three'.
//Function Description
//
// Complete the camelcase function in the editor below.
// camelcase has the following parameter(s):
//
// string s: the string to analyze
// Returns int: the number of words in
//Sample Input : saveChangesInTheEditor
// Sample Output: 5
func main() {
	var input string
	_, _ = fmt.Scanf("%s\n", &input)

	answer := 1
	for _, ch := range input {
		str := string(ch)

		if strings.ToUpper(str) == str {
			answer++
		}

		//min, max := 'A', 'Z'
		//if ch >= min && ch <= max {
		//	// Ch is a capital letter
		//	answer++
		//}
	}

	fmt.Println(answer)
}
