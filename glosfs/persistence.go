package glosfs

import (
	"fmt"
	"os"
	"path/filepath"
)

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
