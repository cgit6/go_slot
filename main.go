package main

import (
	"fmt"
	"os"
	"runtime"
)

// ok

func main() {
	fmt.Println("✅ Hello from Go in Codespaces!")

	// 顯示 Go 版本
	fmt.Println("Go version:", runtime.Version())

	// 顯示目前工作目錄
	if wd, err := os.Getwd(); err == nil {
		fmt.Println("Working directory:", wd)
	}

	// 顯示傳進來的參數（如果有）
	if len(os.Args) > 1 {
		fmt.Println("Args:", os.Args[1:])
	} else {
		fmt.Println("No extra args. Try: go run main.go foo bar")
	}
}
