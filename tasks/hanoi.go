package tasks

func Hanoi(num int) {
	A := num
	B := 0
	C := 0

	// ???

	HanoiSolver(A, B, C)
}

func HanoiSolver(A, B, C int) {
	A++
	C++
	B++

}

func TestHanoiSolver() {
	Hanoi(3)
}
