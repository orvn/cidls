package main

import (
	"fmt"
	"os"

	cid "github.com/ipfs/go-cid"
	multihash "github.com/multiformats/go-multihash"
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

	// Print files and their CIDs
	for _, file := range files {
		fmt.Printf("%s\t", file.Name())
		if !file.IsDir() { // Only process files, not directories
			cidStr, err := computeCID(file.Name())
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			} else {
				fmt.Printf("%s\n", cidStr)
			}
		} else {
			// For directories, just print the name without a CID
			fmt.Println()
		}
	}
}

func computeCID(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	hash, err := multihash.Sum(content, multihash.SHA2_256, -1)
	if err != nil {
		return "", err
	}

	c := cid.NewCidV1(cid.Raw, hash)
	return c.String(), nil
}
