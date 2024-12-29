package utils

import (
	"fmt"
	"os"
	// "path/filepath"

)

func CreateDirAll(path string) error {
	log := Logger
	// Check if the directory already exists.
	_, err := os.Stat(path)
	if err == nil {
		return nil // Directory already exists.
	} else if !os.IsNotExist(err) {
		log.Info("CreateDirAll: Directory already exists", "path", path)
		return nil
	}

	// Create parent directories recursively.
	// parentDir := filepath.Dir(path)
	// if err := CreateDirAll(parentDir); err != nil {
	// 	return fmt.Errorf("error creating parent directory: %w", err)
	// }

	// Attempt to create the final directory.
	err = os.MkdirAll(path, 0755) // Create with permissions 0755 (rwxr-xr-x).
	if err != nil {
		log.Error("CreateDirAll: Failed to create directory", "path", path, "err", err)
		return fmt.Errorf("error creating directory: %s", path)
	}
	log.Info("CreateDirAll: all directories created", "path", path)
	return nil
}

func PathExists(absolutePath string) error {
	_, err := os.Stat(absolutePath)
	if err != nil {
		return err	
	}
	return nil 
}
