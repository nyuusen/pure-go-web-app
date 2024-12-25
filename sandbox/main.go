package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	/*
		ビット演算
	*/
	a := 12 // 1100
	b := 7  // 0111

	// ビット論理積を計算
	and := a & b
	fmt.Printf("AND: %04b\n", and)

	// ビット論理和を計算
	or := a | b
	fmt.Printf("OR: %04b\n", or)

	/*
		ディレクトリ操作関連
	*/
	// ビルド実行時のカレントディレクトリを取得
	dir, _ := os.Getwd()
	fmt.Printf("Getwd: %s\n", dir)
	// このファイルの絶対パスを取得
	absDir, _ := filepath.Abs("./server/main.go")
	fmt.Printf("Abs: %s\n", absDir)
}
