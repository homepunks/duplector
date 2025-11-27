package main

import (
	"log"
	"os"
)

func main() {
	szGroups := make(map[int64][]FileInfo)
	
	log.Println("[INFO] Starting to search for duplicates...")	
	files := ListFiles(os.Getenv("HOME"))

	for _, file := range files {
		szGroups[file.Size] = append(szGroups[file.Size], file)
	}

	log.Printf("[INFO] Found %d files, grouped into %d size groups\n", len(files), len(szGroups))

	dups := make(map[string][]string)

	for size, fileGroup := range szGroups {
		if len(fileGroup) < 2 {
			continue
		}

		log.Printf("[INFO] Checking %d files with size %d bytes\n", len(fileGroup), size)
		
		typeGroups := make(map[string][]FileInfo)
		for _, file := range fileGroup {
			fileType := getFileType(file.Path)
			typeGroups[fileType] = append(typeGroups[fileType], file)
		}

		for fileType, typeGroup := range typeGroups {
			if len(typeGroup) < 2 {
				continue
			}

			log.Printf("[INFO] Checking %d files of type %s\n", len(typeGroup), fileType)
			
			for i := 0; i < len(typeGroup); i++ {
				for j := i + 1; j < len(typeGroup); j++ {
					file1 := typeGroup[i]
					file2 := typeGroup[j]

					const rangeSize = 4096 // 4KB
					sameRanges, err := compareFileRanges(file1.Path, file2.Path, rangeSize)
					if err != nil {
						log.Printf("[WARN] Could not compare ranges for %s and %s: %v", file1.Path, file2.Path, err)
						continue
					}

					if !sameRanges {
						continue
					}

					log.Printf("[INFO] Files passed range check, computing hash: %s and %s", file1.Path, file2.Path)
					
					hash1, err := HashFile(file1.Path)
					if err != nil {
						log.Printf("[ERROR] Could not hash %s: %v", file1.Path, err)
						continue
					}

					hash2, err := HashFile(file2.Path)
					if err != nil {
						log.Printf("[ERROR] Could not hash %s: %v", file2.Path, err)
						continue
					}

					if hash1 == hash2 {
						if _, exists := dups[hash1]; !exists {
							dups[hash1] = []string{file1.Path}
						}
						dups[hash1] = append(dups[hash1], file2.Path)
						log.Printf("[INFO] Duplicate found: %s matches %s", file2.Path, file1.Path)
					}
				}
			}
		}
	}

	if len(dups) == 0 {
		log.Println("[INFO] No duplicates found!")
	} else {
		log.Printf("[INFO] Found %d sets of duplicates:", len(dups))
		for hash, files := range dups {
			log.Printf("[INFO] Hash %s:", hash)
			for i, file := range files {
				if i == 0 {
					log.Printf("  Original: %s", file)
				} else {
					log.Printf("  Duplicate %d: %s", i, file)
				}
			}
		}
	}

	log.Println("[INFO] The search is finished.")	
}
