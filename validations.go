package main

import (
	"bufio"
	"log/slog"
	"strings"
)

func normalizePKGBUILD(content string) string {
	var normalized []string
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "pkgrel=") {
			continue
		}

		checkSum := false

		if strings.HasPrefix(trimmed, "sha256sums") ||
			strings.HasPrefix(trimmed, "sha512sums") ||
			strings.HasPrefix(trimmed, "md5sums") ||
			strings.HasPrefix(trimmed, "b2sums") ||
			strings.HasPrefix(trimmed, "sha1sums") {
			checkSum = true
		}
		for checkSum && !strings.HasSuffix(trimmed, ")") {
			scanner.Scan()
			line = scanner.Text()
			trimmed = strings.TrimSpace(line)
		}
		if checkSum {
			continue
		}

		normalized = append(normalized, line)
	}

	return strings.Join(normalized, "\n")
}

func comparePKGBUILDs(content1, content2 string) bool {
	slog.Info("Removing check sums and pkgrel to compare ...")
	norm1 := normalizePKGBUILD(content1)
	norm2 := normalizePKGBUILD(content2)
	slog.Info("Removed from both files comparing ...")
	return norm1 == norm2
}
