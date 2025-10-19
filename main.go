package main

import (
	"log/slog"
	"os"
)

func main() {
	pkgbuild := NewPkgBuildFromEnv()

	if err := pkgbuild.validate(); err != nil {
		slog.Error("Validation failed", "err", err)
		os.Exit(1)
	}

	if _, err := pkgbuild.generate(); err != nil {
		slog.Error("Generation failed", "err", err)
		os.Exit(1)
	}
	slog.Info("PKGBUILD updated successfully")
}
