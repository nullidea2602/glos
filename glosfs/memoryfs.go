package glosfs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var MemoryFS = make(map[string]string) // In-memory file system
const fsFilename = "memoryfs.dat"      // Persistent storage file
const utilsDir = "utils"

// Save memoryFS to a file
func SaveMemoryFS() {
	file, err := os.Create(fsFilename)
	if err != nil {
		fmt.Println("Error saving memoryFS:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(MemoryFS); err != nil {
		fmt.Println("Error encoding memoryFS:", err)
	}
}

func LoadMemoryFS() {
	// Load existing memoryFS from file
	file, err := os.Open(fsFilename)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&MemoryFS); err != nil {
			fmt.Println("Error decoding memoryFS:", err)
		}
	}

	// Load prewritten Lua scripts if they exist
	loadPrewrittenScripts()
}

func loadPrewrittenScripts() {
	files, err := os.ReadDir(utilsDir)
	if err != nil {
		// If directory does not exist, ignore it
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".lua" {
			filename := file.Name()
			if _, exists := MemoryFS[filename]; !exists { // Avoid overwriting existing files
				content, err := os.ReadFile(filepath.Join(utilsDir, filename))
				if err == nil {
					MemoryFS[filename] = string(content)
					fmt.Printf("Loaded prewritten script: %s\n", filename)
				} else {
					fmt.Printf("Failed to load %s: %v\n", filename, err)
				}
			}
		}
	}
}
