package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
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

	templatePath string
	outputPath   string
	client       AURClient
}

func NewPkgBuild() *PkgBuild {
	return &PkgBuild{
		client: NewAURClient(5 * time.Second),
	}
}
func NewPkgBuildFromEnv() *PkgBuild {
	pkgbuild := NewPkgBuild()

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
	pkgbuild.outputPath = getenv("output_path", "./output/PKGBUILD")
	return pkgbuild
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func (p *PkgBuild) validate() error {
	if p.CliName == "" {
		return fmt.Errorf("CliName is required")
	}
	if len(p.Maintainers) == 0 {
		return fmt.Errorf("At least one Maintainer is required")
	}
	if p.Pkgname == "" {
		return fmt.Errorf("Pkgname is required")
	}
	if p.Version == "" {
		return fmt.Errorf("Version is required")
	}
	if p.Description == "" {
		return fmt.Errorf("Description is required")
	}
	if p.Url == "" {
		return fmt.Errorf("Url is required")
	}
	if len(p.Arch) == 0 {
		return fmt.Errorf("At least one Arch is required")
	}
	if len(p.Licence) == 0 {
		return fmt.Errorf("At least one Licence is required")
	}
	if len(p.Source_x86_64) == 0 {
		return fmt.Errorf("Source_x86_64 is required")
	}
	return nil
}

func (pkgbuild *PkgBuild) generate() (string, error) {
	slog.Info("starting pkgbuild.generate ..")
	PKGBUILD, err := pkgbuild.template()
	if err != nil {
		slog.Error("Failed to template PKGBUILD dumping\n ", "dump", pkgbuild)
		return "", err
	}

	data, err := pkgbuild.client.getAurPackageVersions(pkgbuild.Pkgname)

	if err != nil {
		slog.Error("Failed to fetch package info from AUR")
		return "", err
	}

	if data.version == pkgbuild.Version {
		slog.Warn("AUR version and current version match, this should only be a PKGBUILD update")
		slog.Info("Comparing PKGBUILD to validate")
		aurPKGBUILD, err := pkgbuild.client.fetchPKGBUILD(pkgbuild.Pkgname)
		if err != nil {
			slog.Error("Failed to fetch PKGBUILD from AUR")
			return "", err
		}
		if comparePKGBUILDs(PKGBUILD, aurPKGBUILD) {
			slog.Error("Files match!! should not publish to the AUR without changes to PKGBUILD file or the software Version")
			return "", fmt.Errorf("PKGBUILD already published to AUR")
		}
		slog.Info("New PKGBUILD file increasing pkgrel number to", "pkgrel", data.pkgrel+1)
		pkgbuild.Pkgrel = data.pkgrel + 1

		slog.Info("Templating again")
		PKGBUILD, err = pkgbuild.template()
		if err != nil {
			slog.Error("Failed to template PKGBUILD dumping\n ", "dump", pkgbuild)
			return "", err
		}
	} else {
		slog.Info("New version means new pkgrel")
		pkgbuild.Pkgrel = 1
	}

	if err := writeFile(pkgbuild.outputPath, PKGBUILD); err != nil {
		return "", err
	}
	slog.Info("finished pkgbuild.generate ..")
	return PKGBUILD, nil
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

	if _, err = file.WriteString(content); err != nil {
		return err
	}

	return nil
}

func (pkgbuild PkgBuild) template() (string, error) {
	slog.Info("Templating ...")
	tmpl, err := template.ParseFiles(pkgbuild.templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, pkgbuild); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func comparePKGBUILDs(content1, content2 string) bool {
	slog.Info("Removing check sums and pkgrel to compare ...")
	norm1 := normalizePKGBUILD(content1)
	norm2 := normalizePKGBUILD(content2)
	slog.Info("Removed from both files comparing ...")
	return norm1 == norm2
}
