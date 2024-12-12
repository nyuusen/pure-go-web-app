package main

import "fmt"

func main() {
	a := 12 // 1100
	b := 7  // 0111

	// ビット論理積を計算
	and := a & b
	fmt.Printf("AND: %04b\n", and)

	// ビット論理和を計算
	or := a | b
	fmt.Printf("OR: %04b\n", or)
}
