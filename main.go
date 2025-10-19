package main

import (
	"log/slog"
	"os"
	"strings"
)

func main() {
	pkgbuild := &PkgBuild{}

	pkgbuild.Maintainers = strings.Split(os.Getenv("maintainers"), ",")
	pkgbuild.Contributors = strings.Split(os.Getenv("contributors"), ",")
	pkgbuild.Pkgname = os.Getenv("pkgname")
	pkgbuild.Version = os.Getenv("version")
	pkgbuild.Pkgrel = 1
	pkgbuild.Description = os.Getenv("description")
	pkgbuild.Url = os.Getenv("url")
	pkgbuild.Arch = strings.Split(os.Getenv("arch"), ",")
	pkgbuild.Licence = strings.Split(os.Getenv("licence"), ",")
	pkgbuild.Provides = strings.Split(os.Getenv("provides"), ",")
	pkgbuild.Conflicts = strings.Split(os.Getenv("conflicts"), ",")
	pkgbuild.Source_x86_64 = strings.Split(os.Getenv("source_x86_64"), ",")
	pkgbuild.Source_aarch64 = strings.Split(os.Getenv("source_aarch64"), ",")

	pkgbuild.templatePath = getenv("pkgbuild_template", "./pkgbuild.tmpl")

	if err := pkgbuild.validate(); err != nil {
		slog.Error("Validation failed", "err", err)
		os.Exit(1)
	}

	pkgbuild.generate()
	slog.Info("PKGBUILD updated successfully")
}
func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
