package parser

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
)

func CalculateSHA256(data io.Reader) (string, error) {
	slog.Info("Downloading source for checksum")

	hash := sha256.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
