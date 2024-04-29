package main

import (
	"fmt"
	tasks "golangLearning/tasks"
	"os"
)

func main() {
	idx := -1
	for idx != 0 {
		fmt.Println("\nChoose the task:")
		fmt.Println("\t1 - TextStyleSwitcher")
		fmt.Println("\t2 - MapFromString")
		fmt.Println("\t3 - Calculator")

		fmt.Printf("\nPrint the number of the task to begin (0 to exit):")
		_, err := fmt.Fscan(os.Stdin, &idx)
		if err != nil || idx < 0 {
			fmt.Println("\nError: Incorrect input\n")
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
		default:
			fmt.Printf("\nError: %d is not configured yet\n\n", idx)
		}
	}
}
