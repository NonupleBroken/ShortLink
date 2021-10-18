package util

import (
	"errors"
	"math/rand"
	"os"
	"time"
)

func isPathExists(dirPath string) (bool, error) {
	s, err := os.Stat(dirPath)
	if err == nil {
		if s.IsDir() {
			return true, nil
		} else {
			return false, errors.New(dirPath + " is a file")
		}
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}

func CreateDir(dirPath string) error {
	isExists, err := isPathExists(dirPath)
	if err != nil {
		return err
	}
	if !isExists {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRandomStr(length uint8) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const alphabetLength = len(alphabet)
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = alphabet[rand.Intn(alphabetLength)]
	}
	return string(result)
}