package main

import (
	"log/slog"
	"os"
)

func main() {
	pkgbuild := NewPkgBuildFromEnv()

	slog.Info("Validating", "pkgbuild", pkgbuild)
	if err := validate(*pkgbuild); err != nil {
		slog.Error("Validation failed", "err", err)
		os.Exit(1)
	}
	slog.Info("pkgbuild is valid")

	if _, err := pkgbuild.generate(); err != nil {
		slog.Error("Generation failed", "err", err)
		os.Exit(1)
	}
	slog.Info("PKGBUILD updated successfully")
}
