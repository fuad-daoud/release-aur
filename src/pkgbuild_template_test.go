package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	tests := []struct {
		name             string
		pkg              PkgBuild
		expectedPKGBUILD string
		expectedSRCINFO  string
	}{
		{
			name: "Test all fields",
			pkg: PkgBuild{
				CliName:              "pkg",
				Maintainers:          []string{"Fuad Daoud <aur@fuad-daoud.com>", "Fuad2 Daoud2 <aur2@fuad-daoud.com>"},
				Contributors:         []string{"Someone else <someone@fuad-daoud.com>", "Someone2 else2  <someone2@fuad-daoud.com>"},
				Pkgname:              "pkg-bin",
				Version:              "0.1.4",
				Pkgrel:               1,
				Description:          "Some single line description",
				Url:                  "https://github.com/fuad-daoud/pkg",
				Arch:                 []string{"x86_64", "aarch64"},
				Licence:              []string{"MIT", "OBSD"},
				Provides:             []string{"package-a", "package-b"},
				Conflicts:            []string{"package-c", "package-d"},
				Source_x86_64:        []string{"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64", "LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE", "README::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/README.md"},
				Source_aarch64:       []string{"pkg-bin-0.1.4-aarch_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-aarch_64", "LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE", "README::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/README.md"},
				pkgbuildTemplatePath: "pkgbuild.tmpl",
				srcInfoTemplatePath:  "srcinfo.tmpl",
				client:               NewClient(10*time.Second, time.Second, 5),
			},
			expectedPKGBUILD: "testdata/PKGBUILD",
			expectedSRCINFO:  "testdata/.SRCINFO",
		},
		{
			pkg: PkgBuild{
				CliName:       "pkg",
				Maintainers:   []string{"Fuad Daoud <aur@fuad-daoud.com>", "Fuad2 Daoud2 <aur2@fuad-daoud.com>"},
				Contributors:  []string{"Someone else <someone@fuad-daoud.com>", "Someone2 else2  <someone2@fuad-daoud.com>"},
				Pkgname:       "pkg-bin",
				Version:       "0.1.4",
				Pkgrel:        1,
				Description:   "Some single line description",
				Url:           "https://github.com/fuad-daoud/pkg",
				Arch:          []string{"x86_64", "aarch64"},
				Licence:       []string{"MIT", "OBSD"},
				Provides:      []string{"package-a", "package-b"},
				Conflicts:     []string{"package-c", "package-d"},
				Source_x86_64: []string{"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64", "LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE", "README::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/README.md"},

				pkgbuildTemplatePath: "pkgbuild.tmpl",
				srcInfoTemplatePath:  "srcinfo.tmpl",
				client:               NewClient(10*time.Second, time.Second, 5),
			},
			expectedPKGBUILD: "testdata/PKGBUILD_x86",
			expectedSRCINFO:  "testdata/.SRCINFO_x86",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgbuild, srcinfo, err := tt.pkg.template()

			assert.NoError(t, err)
			expectedPKGBUILD, _ := os.ReadFile(tt.expectedPKGBUILD)
			expectedSRCINFO, _ := os.ReadFile(tt.expectedSRCINFO)
			assert.EqualValuesf(t, string(expectedPKGBUILD), pkgbuild, "Failed Templating")
			assert.EqualValuesf(t, string(expectedSRCINFO), srcinfo, "Failed Templating")
		})
	}
}

func TestTemplate_Errors(t *testing.T) {
	tests := []struct {
		name    string
		pkg     PkgBuild
		wantErr bool
	}{
		{
			name: "invalid template path",
			pkg: PkgBuild{
				pkgbuildTemplatePath: "./nonexistent.tmpl",
			},
			wantErr: true,
		},
		{
			name: "empty template path",
			pkg: PkgBuild{
				pkgbuildTemplatePath: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid template",
			pkg: PkgBuild{
				pkgbuildTemplatePath: "./invalid_template.tmpl",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := tt.pkg.template()
			assert.Error(t, err)
		})
	}
}
