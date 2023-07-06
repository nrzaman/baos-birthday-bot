package util

import (
	"io"
	"os"
)

func Extract(configFile string) []byte {
	// Open the JSON config file
	content, err := os.Open(configFile)
	Check(err)

	defer func(content *os.File) {
		err := content.Close()
		Check(err)
	}(content)

	// Read all contents
	byteResult, err := io.ReadAll(content)
	Check(err)

	return byteResult
}
