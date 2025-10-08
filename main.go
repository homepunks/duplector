package main

import (
	"log"
	"os"
)

func main() {
	checker := make(map[string]string)
	homeDir := os.Getenv("HOME")

	log.Println("[INFO] Starting to search for duplicates...")
	
	files := ListFiles(homeDir)
	for _, file := range files {
		fileHash := HashFile(file)
		if duplicate, found := checker[fileHash]; found {
			log.Printf("[INFO] A duplicate of %s was found: %s\n", file, duplicate)
		}
		checker[fileHash] = file
	}

	log.Println("[INFO] The search is finished.")	
}
