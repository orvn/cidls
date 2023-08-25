package main

import (
	"fmt"
	"os"
	"sync"

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

	// Channel to collect file names and their CIDs
	results := make(chan string, len(files))

	// Use a WaitGroup to wait for all goroutines
	var wg sync.WaitGroup

	// Start a goroutine for each file
	for _, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

			if !file.IsDir() {
				cidStr, err := computeCID(file.Name())
				if err != nil {
					results <- fmt.Sprintf("%s\tERROR: %s", file.Name(), err)
				} else {
					results <- fmt.Sprintf("%s\t%s", file.Name(), cidStr)
				}
			} else {
				results <- fmt.Sprintf("%s\t", file.Name())
			}
		}(file)
	}

	// Close the results channel after all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results as they come in
	for result := range results {
		fmt.Println(result)
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
