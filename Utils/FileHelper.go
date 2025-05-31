package Utils

import (
	"errors"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)

	if nil == err {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func OpenFile(path string) (*os.File, error) {
	if false == FileExists(path) {
		return nil, errors.New("file does not exist")
	}

	return os.Open(path)
}

func ReadFile(path string) ([]byte, error) {
	if false == FileExists(path) {
		return nil, errors.New("file does not exist")
	}

	data, err := os.ReadFile(path)

	return data, err
}
