package main

import (
	"os"
	"testing"

	"github.com/fuad-daoud/release-aur/src/parser"
	"github.com/stretchr/testify/assert"
)

func TestGenerateNewVersion(t *testing.T) {

	pkgbuild := PkgBuild{
		CliName:              "pkgmate",
		Maintainers:          []string{"Fuad Daoud <aur@fuad-daoud.com>"},
		Pkgname:              "pkgmate-bin",
		Version:              "v0.0.0-test-release-aur",
		Pkgrel:               1,
		Description:          "TUI application to manage your dependencies",
		Url:                  "https://github.com/fuad-daoud/pkgmate",
		Arch:                 []string{"x86_64"},
		Licence:              []string{"MIT"},
		Source_x86_64:        []string{"pkgmate-bin-v0.0.0-test-release-aur-x86_64::https://github.com/fuad-daoud/pkgmate/releases/download/v0.0.0-test-release-aur/pkgmate-linux-amd64"},
		pkgbuildTemplatePath: "pkgbuild.tmpl",
		srcInfoTemplatePath:  "srcinfo.tmpl",
		outputPath:           "./output/",
		comparator:           defaultCompareWithRemote,
		checksumCalculator:   parser.DefaultCalculateSources,
	}

	PKGBUILD, err := pkgbuild.generate()
	if err != nil {
		t.Errorf("got an err %v", err)
	}

	expected, err := os.ReadFile("testdata/PKGBUILD_pkgmate")

	if err != nil {
		t.Errorf("got an err %v\n", err)
	}
	assert.EqualValuesf(t, string(expected), PKGBUILD, "Failed Templating")

}

func TestGenerateNewPkgrel(t *testing.T) {
	pkgbuild := PkgBuild{
		CliName:         "pkgmate",
		Maintainers:     []string{"Fuad Daoud <aur@fuad-daoud.com>"},
		Pkgname:         "pkgmate-bin",
		Version:         "0.1.1",
		Pkgrel:          1,
		Description:     "TUI application to manage your dependencies",
		Url:             "https://github.com/fuad-daoud/pkgmate",
		Arch:            []string{"x86_64"},
		Licence:         []string{"MIT"},
		Source_x86_64:   []string{"pkgmate-bin-0.1.1-x86_64::https://github.com/fuad-daoud/pkgmate/releases/download/0.1.1/pkgmate-linux-amd64"},
		Checksum_x86_64: []string{"SKIP"},

		outputPath:           "./output/",
		pkgbuildTemplatePath: "pkgbuild.tmpl",
		srcInfoTemplatePath:  "srcinfo.tmpl",
		comparator:           defaultCompareWithRemote,
		checksumCalculator:   parser.DefaultCalculateSources,
	}

	PKGBUILD, err := pkgbuild.generate()
	if err != nil {
		t.Errorf("got an err %v", err)
	}

	expected, err := os.ReadFile("testdata/PKGBUILD_pkgmate_rel2")

	if err != nil {
		t.Errorf("got an err %v\n", err)
	}
	assert.EqualValuesf(t, string(expected), PKGBUILD, "Failed Templating")

}
func TestGenerate_Errors(t *testing.T) {

	t.Run("template fails", func(t *testing.T) {
		pkg := &PkgBuild{
			CliName:              "test",
			Maintainers:          []string{"User"},
			Pkgname:              "test",
			Version:              "1.0.0",
			Description:          "Test",
			Url:                  "https://example.com",
			Arch:                 []string{"x86_64"},
			Licence:              []string{"MIT"},
			Source_x86_64:        []string{"https://github.com/fuad-daoud/pkgmate/releases/download/0.1.1/pkgmate-linux-amd64"},
			pkgbuildTemplatePath: "./nonexistent.tmpl",
			srcInfoTemplatePath:  "srcinfo.tmpl",
			checksumCalculator:   parser.DefaultCalculateSources,
		}

		_, err := pkg.generate()
		assert.Error(t, err)
	})
	t.Run("PKGBUILDs match - new pkgrel, second templating fails", func(t *testing.T) {

		err := copyFile("pkgbuild.tmpl", "/tmp/pkgbuild.tmpl")
		if err != nil {
			t.Errorf("error in setup, %v", err)
		}

		pkg := &PkgBuild{
			CliName:              "test",
			Maintainers:          []string{"User"},
			Pkgname:              "test",
			Version:              "1.0.0",
			Description:          "Test",
			Url:                  "https://example.com",
			Arch:                 []string{"x86_64"},
			Licence:              []string{"MIT"},
			Source_x86_64:        []string{"https://github.com/fuad-daoud/pkgmate/releases/download/0.1.1/pkgmate-linux-amd64"},
			pkgbuildTemplatePath: "/tmp/pkgbuild.tmpl",
			srcInfoTemplatePath:  "srcinfo.tmpl",
			comparator: func(Client, PkgBuild, string) (int, error) {
				err := copyFile("invalid_template.tmpl", "/tmp/pkgbuild.tmpl")
				if err != nil {
					t.Errorf("error in setup, %v", err)
				}
				return 1, nil

			},

			checksumCalculator: parser.DefaultCalculateSources,
			outputPath:         "/root/",
		}

		_, err = pkg.generate()
		if err == nil {
			t.Errorf("should have retruned an error")
		}
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template")
	})
	t.Run("Geneartion works, fails at writing", func(t *testing.T) {
		pkg := &PkgBuild{
			CliName:              "test",
			Maintainers:          []string{"User"},
			Pkgname:              "test",
			Version:              "1.0.0",
			Description:          "Test",
			Url:                  "https://example.com",
			Arch:                 []string{"x86_64"},
			Licence:              []string{"MIT"},
			Source_x86_64:        []string{"https://github.com/fuad-daoud/pkgmate/releases/download/0.1.1/pkgmate-linux-amd64"},
			Checksum_x86_64:      []string{"SKIP"},
			pkgbuildTemplatePath: "pkgbuild.tmpl",
			srcInfoTemplatePath:  "srcinfo.tmpl",
			outputPath:           "/root/",
			comparator:           func(Client, PkgBuild, string) (int, error) { return -1, nil },
			checksumCalculator:   parser.DefaultCalculateSources,
		}

		_, err := pkg.generate()
		if err == nil {
			t.Errorf("should have retruned an error")
		}
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")
	})
}
