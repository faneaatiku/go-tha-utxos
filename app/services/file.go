package services

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func ReadFile(file string) ([]byte, error) {
	if !FileExists(file) {
		return nil, fmt.Errorf("file %s does not exist", file)
	}

	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func WriteToFileIfNotExists(filename string, data []byte, ignoreExistingFile bool) error {
	if FileExists(filename) {
		if !ignoreExistingFile {
			return fmt.Errorf("file already exists: %s", filename)
		}

		err := os.Remove(filename)
		if err != nil {
			return fmt.Errorf("failed to remove existing file: %s", err)
		}

		log.Infof("removed existing file: %s", filename)
	}

	// open output file
	fo, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error on trying to create file: %w", err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			log.Errorf("error on trying to close file: %v", err)
		}
	}()

	return os.WriteFile(filename, data, os.ModePerm)
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
