package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"sort"
	"sync"

	cid "github.com/ipfs/go-cid"
	multihash "github.com/multiformats/go-multihash"
)

var (
	Version string
	Build   string
)

type LsColors struct {
	DirColor        string
	SymlinkColor    string
	ExecutableColor string
	DotFileColor    string
	CIDColor        string
}

// Default color formatting if no ls_colors variable set in .bashrc or .zshrc
func getLsColors() LsColors {
	defaultColors := LsColors{
		DirColor:        "\033[34m",
		SymlinkColor:    "\033[36m",
		ExecutableColor: "\033[31m",
		DotFileColor:    "\033[37m",
		CIDColor:        "\033[35m",
	}

	lsColorsEnv := os.Getenv("LS_COLORS")
	if lsColorsEnv == "" {
		return defaultColors
	}

	colors := defaultColors
	for _, colorSetting := range split(lsColorsEnv, ":") {
		parts := split(colorSetting, "=")
		if len(parts) != 2 {
			continue
		}

		colorCode := "\033[" + parts[1] + "m"
		switch parts[0] {
		case "di":
			colors.DirColor = colorCode
		case "ln":
			colors.SymlinkColor = colorCode
		case "ex":
			colors.ExecutableColor = colorCode
		case "cid":
			colors.CIDColor = colorCode
		}
	}

	return colors
}

func split(s, sep string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i:i+1] == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// Expand tilde to the full home directory path
func expandTilde(path string) (string, error) {
	if path[:1] == "~" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return usr.HomeDir + path[1:], nil
	}
	return path, nil
}

func main() {
	// Print version and build number
	// fmt.Printf("Version: %s, Build: %s\n", Version, Build)

	// CLI flags
	helpFlag := flag.Bool("h", false, "Display help information")
	versionFlag := flag.Bool("v", false, "Display version information")

	flag.Parse()

	if *helpFlag {
		fmt.Println("Usage: cidls [OPTIONS] [DIRECTORY]")
		fmt.Println("List information about the files in the DIRECTORY (the current directory is the default).")
		fmt.Println("\nOptions:")
		fmt.Println("  -h\tDisplay this help message")
		fmt.Println("  -v\tDisplay version information")
		return
	}

	if *versionFlag {
		fmt.Printf("Version: %s, Build: %s\n\n", Version, Build)
		return
	}

	// Get directory from $1 argument or use the current directory
	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			fmt.Println("Error: unable to get current directory:", err)
			return
		}
	}

	// Get CID version from $1 argument or use default v1
	var cidVersion int = 1
	if len(os.Args) > 2 {
		versionArg := os.Args[2]
		if versionArg == "0" {
			cidVersion = 0
		} else if versionArg != "1" {
			fmt.Println("Error: Invalid CID version. Use 0 or 1.")
			return
		}
	}

	// Expand the tilde ~ to the full home directory path
	dir, err := expandTilde(dir)
	if err != nil {
		fmt.Println("Error expanding tilde:", err)
		return
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

	// Sort the merged list to ensure directories are always at the top
	sort.Slice(sortedFiles, func(i, j int) bool {
		if sortedFiles[i].IsDir() && !sortedFiles[j].IsDir() {
			return true
		}
		if !sortedFiles[i].IsDir() && sortedFiles[j].IsDir() {
			return false
		}
		return sortedFiles[i].Name() < sortedFiles[j].Name()
	})

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
	colors := getLsColors()
	resetColor := "\033[0m"

	// Start a goroutine for each file
	for _, file := range sortedFiles {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

			formattedName := fmt.Sprintf("%-*s", maxNameLen, file.Name())

			if file.Name()[0] == '.' {
				// Handle dotdirectories and dotfiles
				if file.IsDir() {
					results <- fmt.Sprintf("%s%s%s\t", colors.DirColor, formattedName, resetColor)
					return
				}
				cidStr, err := computeCID(dir+string(os.PathSeparator)+file.Name(), cidVersion)
				if err != nil {
					results <- fmt.Sprintf("%s%s%s\tERROR: %s", colors.DotFileColor, formattedName, resetColor, err)
				} else {
					results <- fmt.Sprintf("%s%s%s\t%s%s%s", colors.DotFileColor, formattedName, resetColor, colors.CIDColor, cidStr, resetColor)
				}
				return
			}

			if file.IsDir() {
				results <- fmt.Sprintf("%s%s%s\t", colors.DirColor, formattedName, resetColor)
			} else if (file.Type() & os.ModeSymlink) == os.ModeSymlink {
				results <- fmt.Sprintf("%s%s%s\t", colors.SymlinkColor, formattedName, resetColor)
			} else if (file.Type() & os.ModePerm) == 0111 {
				results <- fmt.Sprintf("%s%s%s\t", colors.ExecutableColor, formattedName, resetColor)
			} else {
				// Check file readability and output any errors
				cidStr, err := computeCID(dir+string(os.PathSeparator)+file.Name(), cidVersion)
				if err != nil {
					if err.Error() == "insufficient permissions" || err.Error() == "no such file" {
						results <- fmt.Sprintf("%s\t\033[31m%s\033[0m", formattedName, err.Error())
					} else {
						results <- fmt.Sprintf(formattedName, err)
					}
				} else {
					results <- fmt.Sprintf("%s\t%s%s%s", formattedName, colors.CIDColor, cidStr, resetColor)
				}
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

func computeCID(filename string, cidVersion int) (string, error) {
	// Attempt to open the file in read-only mode to check permissions
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		if os.IsPermission(err) {
			return "", fmt.Errorf("insufficient permissions")
		}
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no such file")
		}
		return "", err
	}
	file.Close()

	// Now read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	hash, err := multihash.Sum(content, multihash.SHA2_256, -1)
	if err != nil {
		return "", err
	}

	var c cid.Cid
	if cidVersion == 0 {
		c = cid.NewCidV0(hash)
	} else {
		c = cid.NewCidV1(cid.Raw, hash)
	}

	return c.String(), nil
}
