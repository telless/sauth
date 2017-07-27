package utils

import (
	"io/ioutil"
)

func GetFileContent(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	CheckError(err, "open file " + filePath, FatalLogLevel)

	return content
}
