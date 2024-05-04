package main

import (
	"fmt"
	tasks "golangLearning/tasks"
	"os"
)

func main() {
	fmt.Println(tasks.Factorial(5))
	idx := -1
	for idx != 0 {
		fmt.Println("\nChoose the task:")
		fmt.Println("\t1 - TextStyleSwitcher")
		fmt.Println("\t2 - MapFromString")
		fmt.Println("\t3 - PriceInWords")
		fmt.Println("\t4 - Calculate Factorial")
		fmt.Println("\t5 - Calculate Fibonacci")

		fmt.Printf("\nPrint the number of the task to begin (0 to exit):")
		_, err := fmt.Fscan(os.Stdin, &idx)
		if err != nil || idx < 0 {
			fmt.Println("\nError: Incorrect input")
		}
		switch idx {
		case 0:
			fmt.Println("Stopped")
		case 1:
			tasks.TestTextStyleSwitcher()
		case 2:
			tasks.TestMapFromString()
		case 3:
			tasks.TestValueInTextFormat()
		case 4:
			tasks.TestFactorial()
		case 5:
			tasks.TestFibonacci()
		default:
			fmt.Printf("\nError: %d is not configured yet\n\n", idx)
		}
	}
}
