package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestWriteFile_Errors(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		content  string
		setup    func()
		cleanup  func()
	}{
		{
			name:     "invalid path - read-only parent",
			filePath: "/root/test/PKGBUILD",
			content:  "test",
		},
		{
			name:     "permission denied on directory creation",
			filePath: "/test-readonly/subdir/file.txt",
			content:  "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			err := writeFile(tt.filePath, tt.content)
			assert.Error(t, err)
		})
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
