package glosfs

import (
	"encoding/json"
	"fmt"
	"os"
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
