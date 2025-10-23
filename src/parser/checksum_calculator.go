package parser

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func CalculateSHA256(data io.Reader) (string, error) {
	slog.Info("Downloading source for checksum")

	hash := sha256.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

type CalculateSources func(get func(string) ([]byte, error), sources []string) ([]string, error)

func DefaultCalculateSources(get func(string) ([]byte, error), sources []string) ([]string, error) {
	checksums := make([]string, len(sources))

	for i, source := range sources {
		url := source

		if idx := strings.LastIndex(source, "::"); idx != -1 {
			url = source[idx+2:]
		}
		body, err := get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to download source %v: %w", url, err)
		}

		checksum, err := CalculateSHA256(bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("failed to checksum source %d: %w", i, err)
		}

		checksums[i] = checksum
		slog.Info("Calculated checksum", "source", source, "sha256", checksum)
	}

	return checksums, nil
}
