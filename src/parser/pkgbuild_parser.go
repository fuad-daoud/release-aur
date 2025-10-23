package parser

import (
	"bufio"
	"log/slog"
	"strings"
)

func ExtractChecksums(pkgbuildContent string) (map[string][]string, error) {
	checksums := make(map[string][]string)

	scanner := bufio.NewScanner(strings.NewReader(pkgbuildContent))
	var currentArch string
	var inChecksumArray bool

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "sha256sums_") {
			currentArch = line[len("sha256sums_"):strings.Index(line, "=")]
			inChecksumArray = true

			checksums[currentArch] = extractChecksumsFromLine(line)
			if strings.HasSuffix(line, ")") {
				inChecksumArray = false
			}
		} else if inChecksumArray {
			moreChecksum := extractChecksumsFromLine(line)
			checksums[currentArch] = append(checksums[currentArch], moreChecksum...)
			if strings.HasSuffix(line, ")") {
				inChecksumArray = false
			}
		}
	}

	return checksums, scanner.Err()
}

func extractChecksumsFromLine(line string) []string {
	checksums := make([]string, 0)
	if strings.Index(line, "(") != -1 {
		line = strings.TrimPrefix(line, line[:strings.Index(line, "(")+1])
	}
	line = strings.TrimSuffix(line, ")")
	line = strings.TrimSpace(line)

	parts := strings.SplitSeq(line, " ")
	for part := range parts {
		if part != "" {
			checksums = append(checksums, strings.Trim(part, "'"))
		}
	}
	line = strings.Trim(line, "'")
	return checksums
}

func NormalizePKGBUILD(content string) string {
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

func ComparePKGBUILDs(content1, content2 string) bool {
	slog.Info("Removing check sums and pkgrel to compare ...")
	norm1 := NormalizePKGBUILD(content1)
	norm2 := NormalizePKGBUILD(content2)
	slog.Info("Removed from both files comparing ...")
	return norm1 == norm2
}
