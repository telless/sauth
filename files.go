package main

import (
	"io/ioutil"
)

func getFileContent(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	checkError(err, "open file "+filePath, fatalLogLevel)

	return content
}
