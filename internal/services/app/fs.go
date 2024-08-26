package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func (s *Service) getImageFileNames() ([]string, error) {
	fileNames := []string{}

	err := filepath.Walk(s.imagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".jpeg" {

			fileNames = append(fileNames, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk path: %w", err)
	}

	sort.Strings(fileNames)

	return fileNames, nil
}
