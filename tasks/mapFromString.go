package tasks

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Equalize   = '='
	Sum        = '+'
	Separator  = ';'
	Tabulation = "\t"
)

func readExpression(expression string) (string, string, error) {
	parts := strings.Split(expression, string(Equalize))
	if len(parts) != 2 {
		return "", "", fmt.Errorf("incorrect expression")
	}
	return parts[0], parts[1], nil
}

func MapFromString(text string) (map[string]string, error) {
	text = strings.ReplaceAll(text, string(Space), "")
	text = strings.ReplaceAll(text, string(Tabulation), "")

	inputList := strings.Split(text, string(Separator))
	varMap := make(map[string]string)

	for i, input := range inputList {
		key, value, err := readExpression(input)
		if err != nil {
			err = fmt.Errorf("error while processing %d argument, %w", i, err)
			return varMap, err
		}
		varMap[key] = value
	}

	return varMap, nil
}

func TestMapFromString() {
	fmt.Println("\nThis utility makes a map from a string dataset, for example from:")
	testString := "a= 2; b=sdfadg; c= 1 4. 123; d =fj2, a34, 444"
	fmt.Println(testString)
	fmt.Println("will be done the next map:")
	fmt.Println(MapFromString(testString))

	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')

	fmt.Println(MapFromString(str))
}
