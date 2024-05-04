package tasks

import (
	"fmt"
	"os"
)

func Factorial(v int) int {
	if v == 0 {
		return 1
	}
	return v * Factorial(v-1)
}

func TestFactorial() {
	fmt.Printf("\nPrint the number to get it's factorial : ")
	var number int
	n, err := fmt.Fscan(os.Stdin, &number)
	if err != nil || n > 1 {
		fmt.Printf("\nIncorrect input\n")
	}
	fmt.Printf("\nResult: %d\n", Factorial(number))
}
