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
		_ = content.Close() // Best effort close, error handled by Check(err) above if needed
	}(content)

	// Read all contents
	byteResult, err := io.ReadAll(content)
	Check(err)

	return byteResult
}
