package parser

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateSHA256(t *testing.T) {
	tests := []struct {
		name     string
		input    io.Reader
		expected string
		wantErr  bool
	}{
		{
			name:     "known content",
			input:    strings.NewReader("test content"),
			expected: "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72",
			wantErr:  false,
		},
		{
			name:     "empty content",
			input:    strings.NewReader(""),
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr:  false,
		},
		{
			name:     "SKIP",
			input:    strings.NewReader("SKIP"),
			expected: "6ad446059a8bb8d8722e1420110f2784368a2ad5a56384d516c602530e5af256",
			wantErr:  false,
		},
		{
			name:     "multiline content",
			input:    strings.NewReader("line1\nline2\nline3"),
			expected: "6bb6a5ad9b9c43a7cb535e636578716b64ac42edea814a4cad102ba404946837",
			wantErr:  false,
		},
		{
			name:    "reader error",
			input:   &errorReader{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateSHA256(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to calculate checksum")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock read error")
}

func TestCalculateSHA256_LargeFile(t *testing.T) {
	// Test with 10MB of data
	largeData := bytes.Repeat([]byte("a"), 10*1024*1024)
	reader := bytes.NewReader(largeData)

	hash, err := CalculateSHA256(reader)

	assert.NoError(t, err)
	assert.Len(t, hash, 64)
}

func TestCalculateSHA256_Deterministic(t *testing.T) {
	input := "deterministic test"

	hash1, err1 := CalculateSHA256(strings.NewReader(input))
	hash2, err2 := CalculateSHA256(strings.NewReader(input))

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, hash1, hash2)
}
