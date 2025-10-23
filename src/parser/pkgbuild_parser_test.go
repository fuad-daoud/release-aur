package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractChecksums(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string][]string
	}{
		{
			name: "single SKIP checksum",
			input: `pkgname=test
sha256sums_x86_64=('SKIP')`,
			expected: map[string][]string{
				"x86_64": {"SKIP"},
			},
		},
		{
			name:  "multiple checksums single line",
			input: `sha256sums_x86_64=('abc123' 'def456' 'ghi789')`,
			expected: map[string][]string{
				"x86_64": {"abc123", "def456", "ghi789"},
			},
		},
		{
			name: "multi-line checksums",
			input: `sha256sums_x86_64=('7f6d936dae7da64b45acdadef654863f2d47866660c6cf821430707bbf63c4cd'
                   'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4'
                   '8320695d303094310734f5df7e96722a03d7948b076d52849a7b014006aff793')`,
			expected: map[string][]string{
				"x86_64": {
					"7f6d936dae7da64b45acdadef654863f2d47866660c6cf821430707bbf63c4cd",
					"c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
					"8320695d303094310734f5df7e96722a03d7948b076d52849a7b014006aff793",
				},
			},
		},
		{
			name: "both x86_64 and aarch64",
			input: `sha256sums_x86_64=('abc123' 'def456')
sha256sums_aarch64=('ghi789' 'jkl012')`,
			expected: map[string][]string{
				"x86_64":  {"abc123", "def456"},
				"aarch64": {"ghi789", "jkl012"},
			},
		},
		{
			name: "real PKGBUILD from testdata",
			input: `pkgname=pkg-bin
pkgver=0.1.4
pkgrel=1
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/..."
"LICENSE::https://raw.githubusercontent.com/..."
"README::https://raw.githubusercontent.com/..."
)
sha256sums_x86_64=('SKIP')
sha256sums_aarch64=('SKIP')`,
			expected: map[string][]string{
				"x86_64":  {"SKIP"},
				"aarch64": {"SKIP"},
			},
		},
		{
			name:     "no checksums",
			input:    `pkgname=test\npkgver=1.0.0`,
			expected: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractChecksums(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractChecksumsFromLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "with prefix and suffix",
			input:    "sha256sums_x86_64=('SKIP')",
			expected: []string{"SKIP"},
		},
		{
			name:     "continuation line with quotes",
			input:    "                   'abc123'",
			expected: []string{"abc123"},
		},
		{
			name:     "closing parenthesis only",
			input:    "                   'def456')",
			expected: []string{"def456"},
		},
		{
			name:     "middle line in array",
			input:    "'7f6d936dae7da64b45acdadef654863f2d47866660c6cf821430707bbf63c4cd'",
			expected: []string{"7f6d936dae7da64b45acdadef654863f2d47866660c6cf821430707bbf63c4cd"},
		},

		{
			name:     "multiples in the same line",
			input:    "'abc' 'sd' 'fd')",
			expected: []string{"abc", "sd", "fd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractChecksumsFromLine(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComparePKGBUILDs(t *testing.T) {
	pkgbuild1 := `pkgname=test
pkgver=1.0.0
pkgrel=1
sha256sums=('abc123')`

	pkgbuild2 := `pkgname=test
pkgver=1.0.0
pkgrel=5
sha256sums=('different456')`

	if !ComparePKGBUILDs(pkgbuild1, pkgbuild2) {
		t.Error("Expected PKGBUILDs to be equal when ignoring checksums and pkgrel")
	}

	pkgbuild3 := `pkgname=different
pkgver=1.0.0
pkgrel=9`

	if ComparePKGBUILDs(pkgbuild1, pkgbuild3) {
		t.Error("Expected PKGBUILDs to be different")
	}
}
func TestNormalizePKGBUILD(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "removes checksums and pkgrel",
			input: `pkgname=test
pkgver=1.0.0
pkgrel=1
sha256sums=('abc123')
sha256sums_x86_64=('def456')`,
			expected: `pkgname=test
pkgver=1.0.0`,
		},
		{
			name: "removes multi-line checksum arrays",
			input: `pkgname=test
pkgver=1.0.0
pkgrel=2
sha256sums_x86_64=('7f6d936dae7da64b45acdadef654863f2d47866660c6cf821430707bbf63c4cd'
                   'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4'
                   '8320695d303094310734f5df7e96722a03d7948b076d52849a7b014006aff793')
sha256sums_aarch64=('3bc82df11d552c8134436fd0d752cda0f551932d9e80f038110f0c3ec0e39232'
                    'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4'
                    '8320695d303094310734f5df7e96722a03d7948b076d52849a7b014006aff793')
arch=('x86_64')`,
			expected: `pkgname=test
pkgver=1.0.0
arch=('x86_64')`,
		},
		{
			name: "keeps source arrays",
			input: `pkgname=test
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
pkgrel=1
sha256sums=('SKIP')`,
			expected: `pkgname=test
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)`,
		},
		{
			name: "handles SKIP checksums",
			input: `pkgname=test
sha256sums=('SKIP')
sha256sums_x86_64=('SKIP')`,
			expected: `pkgname=test`,
		},
		{
			name: "full PKGBUILD example",
			input: `# Maintainer: Test User
pkgname=prayers-bin
pkgver=0.1.4
pkgrel=1
pkgdesc="Test package"
arch=('x86_64' 'aarch64')
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
source_aarch64=(
"pkg-bin-0.1.4-aarch_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-aarch_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
sha256sums_x86_64=('7f6d936dae7da64b45acdadef654863f2d47866660c6cf821430707bbf63c4cd'
                   'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4')
sha256sums_aarch64=('3bc82df11d552c8134436fd0d752cda0f551932d9e80f038110f0c3ec0e39232'
                    'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4')
package() {
    install -Dm755 "$srcdir/pkg-bin" "$pkgdir/usr/bin/prayers"
}`,
			expected: `# Maintainer: Test User
pkgname=prayers-bin
pkgver=0.1.4
pkgdesc="Test package"
arch=('x86_64' 'aarch64')
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
source_aarch64=(
"pkg-bin-0.1.4-aarch_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-aarch_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
package() {
    install -Dm755 "$srcdir/pkg-bin" "$pkgdir/usr/bin/prayers"
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePKGBUILD(tt.input)
			assert.EqualValuesf(t, tt.expected, result, "Failed Normalizing")
		})
	}
}
