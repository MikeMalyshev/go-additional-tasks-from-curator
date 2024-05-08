package tasks

import (
	"fmt"
	"os"
)

type rod struct {
	name string
	num  int
}

func TestHanoiSolver() {
	fmt.Printf("\nPrint the start number of disks on the first rod (0 to exit):")

	var idx int
	_, err := fmt.Fscan(os.Stdin, &idx)
	if err != nil || idx < 0 {
		fmt.Println("\nError: Incorrect input")
	}
	if idx == 0 {
		return
	}
	if idx > 30 {
		fmt.Println("Too many disks to begin")
		return
	}
	fmt.Println("Steps:")
	Hanoi(idx)
}

func Hanoi(num int) {
	A := rod{"'A'", num}
	B := rod{"'B'", 0}
	C := rod{"'C'", 0}

	hanoiSolver(&A, &B, &C, num)
}

func hanoiSolver(src, dst, tmp *rod, disksToMove int) {
	if disksToMove > 0 {
		hanoiSolver(src, tmp, dst, disksToMove-1)
		moveOne(src, dst)
		hanoiSolver(tmp, dst, src, disksToMove-1)
	}
}

func moveOne(src, dst *rod) {
	fmt.Printf("\tmove from %s(n=%d) to %s(n=%d)", src.name, src.num, dst.name, dst.num)
	src.num--
	dst.num++
	fmt.Printf(" -> %s(n=%d); %s(n=%d)\n", src.name, src.num, dst.name, dst.num)
}
