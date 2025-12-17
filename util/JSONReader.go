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
		if closeErr := content.Close(); closeErr != nil {
			Check(closeErr)
		}
	}(content)

	// Read all contents
	byteResult, err := io.ReadAll(content)
	Check(err)

	return byteResult
}
