package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePKGBUILD(t *testing.T) {
	for _, tt := range normalize_tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePKGBUILD(tt.input)
			assert.EqualValuesf(t, tt.expected, result, "Failed Normalizing")
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

	if !comparePKGBUILDs(pkgbuild1, pkgbuild2) {
		t.Error("Expected PKGBUILDs to be equal when ignoring checksums and pkgrel")
	}

	pkgbuild3 := `pkgname=different
pkgver=1.0.0
pkgrel=9`

	if comparePKGBUILDs(pkgbuild1, pkgbuild3) {
		t.Error("Expected PKGBUILDs to be different")
	}
}
