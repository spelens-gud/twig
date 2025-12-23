package kfile

import (
	"io"
	"os"
)

// LoadFileByPath function 加载文件.
func LoadFileByPath(path string) (data []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	//nolint: errcheck
	defer file.Close()

	data, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
