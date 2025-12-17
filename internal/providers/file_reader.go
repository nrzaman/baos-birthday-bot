package providers

import (
	"io"
	"os"
)

// RealFileReader implements FileReader using actual file system
type RealFileReader struct{}

func (r *RealFileReader) ReadFile(path string) ([]byte, error) {
	content, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer content.Close()

	byteResult, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}

	return byteResult, nil
}
