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

	"github.com/fuad-daoud/release-aur/src/parser"
)

type PkgBuild struct {
	CliName          string
	Maintainers      []string
	Contributors     []string
	Pkgname          string
	Version          string
	Pkgrel           int
	Description      string
	Url              string
	Arch             []string
	Licence          []string
	Provides         []string
	Conflicts        []string
	Source_x86_64    []string
	Checksum_x86_64  []string
	Source_aarch64   []string
	Checksum_aarch64 []string

	pkgbuildTemplatePath string
	srcInfoTemplatePath  string
	outputPath           string
	comparator           compareWithRemote
	checksumCalculator   parser.CalculateSources
}

func NewPkgBuild() *PkgBuild {
	return &PkgBuild{
		comparator:         defaultCompareWithRemote,
		checksumCalculator: parser.DefaultCalculateSources,
	}
}
func NewPkgBuildFromEnv() *PkgBuild {
	pkgbuild := NewPkgBuild()

	pkgbuild.Maintainers = strings.Split(os.Getenv("maintainers"), ",")
	pkgbuild.Contributors = strings.Split(os.Getenv("contributors"), ",")
	if len(pkgbuild.Contributors) == 1 && pkgbuild.Contributors[0] == "" {
		pkgbuild.Contributors = []string{}
	}
	pkgbuild.CliName = os.Getenv("cli_name")
	pkgbuild.Pkgname = os.Getenv("pkgname")
	pkgbuild.Version = os.Getenv("version")
	if strings.HasPrefix(pkgbuild.Version, "v") && strings.Count(pkgbuild.Version, ".") >= 2 {
		slog.Info("Version starts with 'v' and contains two or more '.', so removing the 'v'", "count", strings.Count(pkgbuild.Version, "."))
		pkgbuild.Version = pkgbuild.Version[1:]
	}
	pkgbuild.Pkgrel = 1
	pkgbuild.Description = os.Getenv("description")
	pkgbuild.Url = os.Getenv("url")
	pkgbuild.Arch = strings.Split(os.Getenv("arch"), ",")
	pkgbuild.Licence = strings.Split(os.Getenv("licence"), ",")
	pkgbuild.Provides = strings.Split(os.Getenv("provides"), ",")
	if len(pkgbuild.Provides) == 1 && pkgbuild.Provides[0] == "" {
		pkgbuild.Provides = []string{}
	}
	pkgbuild.Conflicts = strings.Split(os.Getenv("conflicts"), ",")
	if len(pkgbuild.Conflicts) == 1 && pkgbuild.Conflicts[0] == "" {
		pkgbuild.Conflicts = []string{}
	}
	pkgbuild.Source_x86_64 = strings.Split(os.Getenv("source_x86_64"), ",")
	pkgbuild.Source_aarch64 = strings.Split(os.Getenv("source_aarch64"), ",")

	if len(pkgbuild.Source_aarch64) == 1 && pkgbuild.Source_aarch64[0] == "" {
		pkgbuild.Source_aarch64 = []string{}
	}

	pkgbuild.pkgbuildTemplatePath = getenv("pkgbuild_template", "./pkgbuild.tmpl")
	pkgbuild.srcInfoTemplatePath = getenv("srcinfo_template", "./srcinfo.tmpl")
	pkgbuild.outputPath = getenv("output_path", "./output/")
	return pkgbuild
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func (pkgbuild *PkgBuild) generate() (string, error) {
	slog.Info("starting pkgbuild.generate ..")

	client := NewClient(time.Second*30, time.Second*5, 5)
	var err error
	if pkgbuild.Checksum_x86_64, err = pkgbuild.checksumCalculator(client.Get, pkgbuild.Source_x86_64); err != nil {
		return "", err
	}

	if pkgbuild.Checksum_aarch64, err = pkgbuild.checksumCalculator(client.Get, pkgbuild.Source_aarch64); err != nil {
		return "", err
	}

	PKGBUILD, SRCINFO, err := pkgbuild.template()
	if err != nil {
		slog.Error("Failed to template PKGBUILD dumping\n ", "dump", pkgbuild)
		return "", err
	}

	remotePkgrel, err := pkgbuild.comparator(client, *pkgbuild, PKGBUILD)
	if err != nil {
		return "", err
	}
	if remotePkgrel != -1 {
		pkgbuild.Pkgrel = remotePkgrel + 1
		slog.Info("New PKGBUILD file increasing pkgrel number to", "pkgrel", pkgbuild.Pkgrel)

		slog.Info("Templating again")
		PKGBUILD, SRCINFO, err = pkgbuild.template()
		if err != nil {
			slog.Error("Failed to template PKGBUILD dumping\n ", "dump", pkgbuild)
			return "", err
		}
	}
	if err := writeFile(pkgbuild.outputPath+"PKGBUILD", PKGBUILD); err != nil {
		return "", err
	}
	slog.Info("Wrote PKGBUILD")

	if err := writeFile(pkgbuild.outputPath+".SRCINFO", SRCINFO); err != nil {
		return "", err
	}

	slog.Info("Wrote .SRCINFO")

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

func (pkgbuild PkgBuild) template() (string, string, error) {
	slog.Info("Templating ...")
	tmpl := template.New("pkgbuild").Funcs(template.FuncMap{
		"join_quoted": func(items []string, sep string) string {
			quoted := make([]string, len(items))
			for i, item := range items {
				quoted[i] = "'" + item + "'"
			}
			return strings.Join(quoted, sep)
		},
	})
	tmpl, err := tmpl.ParseFiles(pkgbuild.pkgbuildTemplatePath, pkgbuild.srcInfoTemplatePath)
	if err != nil {
		return "", "", err
	}

	var pkgbuildBuf bytes.Buffer
	templateName := filepath.Base(pkgbuild.pkgbuildTemplatePath)
	if err := tmpl.ExecuteTemplate(&pkgbuildBuf, templateName, pkgbuild); err != nil {
		return "", "", err
	}

	var srcinfoBuf bytes.Buffer
	templateName = filepath.Base(pkgbuild.srcInfoTemplatePath)
	if err := tmpl.ExecuteTemplate(&srcinfoBuf, templateName, pkgbuild); err != nil {
		return "", "", err
	}

	return pkgbuildBuf.String(), srcinfoBuf.String(), nil
}
