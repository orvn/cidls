package main

import (
	"fmt"
	"os"
	"sort"
	"sync"

	cid "github.com/ipfs/go-cid"
	multihash "github.com/multiformats/go-multihash"
)

func main() {
	// Get directory from $1 argument or use the current directory
	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}
	}

	// Check if directory exists and if it's readable
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Error: Directory %s does not exist.\n", dir)
		return
	} else if os.IsPermission(err) {
		fmt.Printf("Error: Permission denied for directory %s, try using sudo?\n", dir)
		return
	}

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:\n", err)
		return
	}

	// Separate directories and files
	var dirs, regularFiles []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		} else {
			regularFiles = append(regularFiles, file)
		}
	}

	// Sort directories and files
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	sort.Slice(regularFiles, func(i, j int) bool { return regularFiles[i].Name() < regularFiles[j].Name() })

	// Merge directories and files
	sortedFiles := append(dirs, regularFiles...)

	// Compute the maximum filename length
	maxNameLen := 0
	for _, file := range sortedFiles {
		if len(file.Name()) > maxNameLen {
			maxNameLen = len(file.Name())
		}
	}

	// Channel to collect file names and their CIDs
	results := make(chan string, len(sortedFiles))

	// Use a WaitGroup to wait for all goroutines
	var wg sync.WaitGroup

	// Start a goroutine for each file
	for _, file := range sortedFiles {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

			formattedName := fmt.Sprintf("%-*s", maxNameLen, file.Name())
			if !file.IsDir() {
				cidStr, err := computeCID(dir + string(os.PathSeparator) + file.Name())
				if err != nil {
					results <- fmt.Sprintf("%s\tERROR: %s", formattedName, err)
				} else {
					results <- fmt.Sprintf("%s\t%s", formattedName, cidStr)
				}
			} else {
				results <- fmt.Sprintf("%s\t", formattedName)
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
