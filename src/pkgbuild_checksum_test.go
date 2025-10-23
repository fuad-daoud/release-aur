package main

import (
	"fmt"
	"testing"

	"github.com/fuad-daoud/release-aur/src/parser"
	"github.com/stretchr/testify/assert"
)

// func TestCalculateForSourcesReal(t *testing.T) {
// 	sources := []string{"pkgmate-linux-amd64::https://github.com/fuad-daoud/pkgmate/releases/download/0.1.1/pkgmate-linux-amd64"}
//
// 	pkgbuild := &PkgBuild{
// 		// client: NewClient(5*time.Second, 1*time.Second, 1),
// 	}
//
// 	result, err := pkgbuild.calculateForSources(sources)
//
// 	assert.NoError(t, err)
// 	// checksum from github
// 	assert.Equal(t, []string{"ba62160a8721ea41c112adc3cd369e1c7abb9d1c03d2bd89d13740b420cc1cc6"}, result)
// }

func TestCalculateForSources(t *testing.T) {
	tests := []struct {
		name        string
		sources     []string
		clientMock  func(string) ([]byte, error)
		wantErr     bool
		errContains string
		expected    []string
	}{

		{
			name:    "single source with known content",
			sources: []string{"https://example.com/file.tar.gz"},
			clientMock: func(string) ([]byte, error) {
				return []byte("test content"), nil
			},
			expected: []string{"6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"},
		},
		{
			name: "multiple sources",
			sources: []string{
				"file1::https://example.com/file1",
				"file2::https://example.com/file2",
			},
			clientMock: func(url string) ([]byte, error) {
				switch url {
				case "https://example.com/file1":
					return []byte("content1"), nil
				case "https://example.com/file2":
					return []byte("content2"), nil
				}
				return []byte("should not return"), nil
			},
			expected: []string{
				"d0b425e00e15a0d36b9b361f02bab63563aed6cb4665083905386c55d5b679fa",
				"dab741b6289e7dccc1ed42330cae1accc2b755ce8079c2cd5d4b5366c9f769a6",
			},
		},
		{
			name:    "source with filename prefix",
			sources: []string{"myfile-1.0.0::https://example.com/release.tar.gz"},
			clientMock: func(string) ([]byte, error) {
				return []byte("release"), nil
			},
			expected: []string{"a4d451ec23463726f72c43d64c710968f6b602cd653b4de8adee1b556240a829"},
		},
		{
			name:    "download fails",
			sources: []string{"https://example.com/notfound"},
			clientMock: func(string) ([]byte, error) {
				return []byte{}, fmt.Errorf("404 Not Found")
			},
			wantErr:     true,
			errContains: "failed to download",
		},
		{
			name:    "empty source list",
			sources: []string{},
			clientMock: func(string) ([]byte, error) {
				return []byte("shouldn't be called"), nil
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.DefaultCalculateSources(tt.clientMock, tt.sources)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
