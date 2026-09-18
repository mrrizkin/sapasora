package storage

import (
	"os"
	"path/filepath"
)

type Local struct {
	config *LocalConfig
}

func NewLocal(config *StorageConfig) *Local {
	return &Local{
		config: &config.Local,
	}
}

func (l *Local) Save(path string, data []byte) error {
	localPath := l.config.Path
	if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(localPath, path)

	if err := os.WriteFile(filePath, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (l *Local) Read(path string) ([]byte, error) {
	localPath := l.config.Path
	filePath := filepath.Join(localPath, path)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (l *Local) Delete(path string) error {
	localPath := l.config.Path
	filePath := filepath.Join(localPath, path)

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func (l *Local) Exists(path string) bool {
	localPath := l.config.Path
	filePath := filepath.Join(localPath, path)

	_, err := os.Stat(filePath)
	return err == nil
}

func (l *Local) List(path string) ([]string, error) {
	localPath := l.config.Path
	files, err := os.ReadDir(localPath)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePaths = append(filePaths, file.Name())
	}

	return filePaths, nil
}

func (l *Local) Move(from, to string) error {
	localPath := l.config.Path
	fromPath := filepath.Join(localPath, from)
	toPath := filepath.Join(localPath, to)

	if err := os.Rename(fromPath, toPath); err != nil {
		return err
	}

	return nil
}
