package utils

import (
	"io/ioutil"
)

func GetFileContent(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	CheckError(err)

	return content
}
