package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type PkgBuild struct {
	CliName        string
	Maintainers    []string
	Contributors   []string
	Pkgname        string
	Version        string
	Pkgrel         int
	Description    string
	Url            string
	Arch           []string
	Licence        []string
	Provides       []string
	Conflicts      []string
	Source_x86_64  []string
	Source_aarch64 []string
}

func main() {
	pkgbuild := PkgBuild{}

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

	PKGBUILD, err := templatePKGBUILD(pkgbuild)
	if err != nil {
		slog.Error("Failed to template PKGBUILD dumping\n ", "dump", pkgbuild)
		os.Exit(1)
	}

	data, err := getAurPackageVersions(pkgbuild.Pkgname)

	if err != nil {
		slog.Error("Failed to fetch package info from AUR", "err", err)
		os.Exit(1)
	}

	if data.version == pkgbuild.Version {
		slog.Warn("AUR version and current version match, this should only be a PKGBUILD update")
		slog.Info("Comparing PKGBUILD to validate")
		aurPKGBUILD, err := fetchPKGBUILD(pkgbuild.Pkgname)
		if err != nil {
			slog.Error("Failed to fetch PKGBUILD from AUR", "err", err)
			os.Exit(1)
		}
		if comparePKGBUILDs(PKGBUILD, aurPKGBUILD) {
			slog.Error("Files match!! should not publish to the AUR without changes to PKGBUILD file or the software Version")
			os.Exit(1)
		}
		slog.Info("New PKGBUILD file increasing pkgrel number to", "pkgrel", data.pkgrel+1)
		pkgbuild.Pkgrel = data.pkgrel + 1

		slog.Info("Templating again")
		PKGBUILD, err = templatePKGBUILD(pkgbuild)
		if err != nil {
			slog.Error("Failed to template PKGBUILD dumping\n ", "dump", pkgbuild)
			os.Exit(1)
		}
	} else {
		slog.Info("New version means new pkgrel")
		pkgbuild.Pkgrel = 1
	}
	writeFile("./output/PKGBUILD", PKGBUILD)
	slog.Info("PKGBUILD updated successfully")
}

func writeFile(filePath string, content string) error {
	outputDir := filepath.Dir(filePath)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func templatePKGBUILD(pkgbuild PkgBuild) (string, error) {
	slog.Info("Templating")
	templateFile := "./pkgbuild.tmpl"
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pkgbuild)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
