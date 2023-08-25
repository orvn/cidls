package main

import (
	"fmt"
	"os"
)

func main() {
	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Print files and directories
	for _, file := range files {
		fmt.Printf("%s\t", file.Name())
	}
	fmt.Println()
}
