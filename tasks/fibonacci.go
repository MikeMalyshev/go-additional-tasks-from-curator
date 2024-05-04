package tasks

import (
	"fmt"
	"os"
)

func Fibonacci(n int) int {
	if n == 1 {
		return 1
	} else if n == 0 {
		return 0
	}
	//  7: 0,1,1,2,3,5,8,13
	return Fibonacci(n-1) + Fibonacci(n-2)
}

func TestFibonacci() {
	fmt.Printf("\nPrint the number to get n-number Fibonacci : ")
	var number int
	n, err := fmt.Fscan(os.Stdin, &number)
	if err != nil || n > 1 {
		fmt.Printf("\nIncorrect input\n")
	}
	fmt.Printf("\nResult: %d\n", Fibonacci(number))
}
