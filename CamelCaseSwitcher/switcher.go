package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const (
	UnderLine = '_'
	Space     = ' '

	UnderLinedStyle = 1
	CamelStyle      = 2
)

func switchToCamelCase(word string) string {
	var output string

	// Setting "previousRune" to UnderLine to convert first literal to uppercase
	previousRune := UnderLine
	for _, s := range word {
		if s != UnderLine {
			if previousRune != UnderLine {
				output += string(s)
			} else {
				output += string(unicode.ToUpper(s))
			}
		}
		previousRune = s
	}

	return output
}

func switchToUnderScore(word string) string {
	var output string

	// Setting "previousWasUpper" to true to convert first literal to lowercase
	previousWasUpper := true
	for _, s := range word {
		if unicode.IsUpper(s) {

			// Вo not add an underline if there are several uppercase letters in a row
			if !previousWasUpper {
				output += string(UnderLine)
			}
			previousWasUpper = true
		} else {
			previousWasUpper = false
		}
		output += string(unicode.ToLower(s))
	}
	return output
}

func checkWordWritingStyle(word string) int {
	if strings.Contains(word, string(UnderLine)) {
		return UnderLinedStyle
	}
	return CamelStyle
}

func processWord(word string) string {
	if checkWordWritingStyle(word) == CamelStyle {
		return switchToUnderScore(word)
	}
	return switchToCamelCase(word)
}

func processString(text string) string {
	inputWordList := strings.Split(text, string(Space))

	var output string
	for _, word := range inputWordList {
		switchedWord := processWord(word)
		output += (switchedWord + " ")
	}
	output, _ = strings.CutSuffix(output, string(Space))

	return output
}

func main() {
	fmt.Println("Utility converts CamelStyle text to underline_style, for example")
	testString := "\"NothingToDoWithIt At ALL, i'm shure! But_when_i write so, something_can_happen\""
	fmt.Print(testString)
	fmt.Print(" will be converted to ")
	fmt.Println(processString(testString))

	fmt.Println("\nWaiting your text:")

	reader := bufio.NewReader(os.Stdin)

	var text, str string
	str, _ = reader.ReadString(';')
	text += str

	fmt.Println(processString(text))
}
