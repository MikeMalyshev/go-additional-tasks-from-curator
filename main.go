package main

import (
	"fmt"
	tasks "golangLearning/tasks"
	"os"
)

func main() {
	idx := -1
	for idx != 0 {
		fmt.Println("Choose the task. Print the number of the task to begin (0 to exit):")
		fmt.Println("1 - TextStyleSwitcher")
		fmt.Println("2 - Calculator")
		fmt.Fscan(os.Stdin, &idx)
		switch idx {
		case 0:
			fmt.Println("Stopped")
		case 1:
			tasks.TextStyleSwitcher()
		default:
			fmt.Printf("\n0Error: %d is not configured yet\n", idx)
		}
	}
}
