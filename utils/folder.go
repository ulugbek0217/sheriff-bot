package utils

import (
	"log"
	"os"
)

func FolderExists(path string) bool {
	status, err := os.Stat(path)
	if err != nil {
		log.Printf("error checking folder: %v\n", err)
		return false
	}
	return status.IsDir()

}
