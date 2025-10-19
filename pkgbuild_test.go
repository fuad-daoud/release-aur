package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	data := PkgBuild{
		CliName:        "pkg",
		Maintainers:    []string{"Fuad Daoud <aur@fuad-daoud.com>", "Fuad2 Daoud2 <aur2@fuad-daoud.com>"},
		Contributors:   []string{"Someone else <someone@fuad-daoud.com>", "Someone2 else2  <someone2@fuad-daoud.com>"},
		Pkgname:        "pkg-bin",
		Version:        "0.1.4",
		Pkgrel:         1,
		Description:    "Some single line description",
		Url:            "https://github.com/fuad-daoud/pkg",
		Arch:           []string{"x86_64", "aarch_64"},
		Licence:        []string{"MIT", "OBSD"},
		Provides:       []string{"package-a", "package-b"},
		Conflicts:      []string{"package-c", "package-d"},
		Source_x86_64:  []string{"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64", "LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE", "README::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/README.md"},
		Source_aarch64: []string{"pkg-bin-0.1.4-aarch_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-aarch_64", "LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE", "README::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/README.md"},

		templatePath: "pkgbuild.tmpl",
	}

	buf, err := data.template()

	if err != nil {
		fmt.Printf("got an err %v\n", err)
		t.FailNow()
	}

	expected, err := os.ReadFile("PKGBUILD")

	if err != nil {
		fmt.Printf("got an err %v\n", err)
		t.FailNow()
	}
	assert.EqualValuesf(t, string(expected), buf, "Failed Templating")
}

func TestX86Only(t *testing.T) {
	data := PkgBuild{
		CliName:       "pkg",
		Maintainers:   []string{"Fuad Daoud <aur@fuad-daoud.com>", "Fuad2 Daoud2 <aur2@fuad-daoud.com>"},
		Contributors:  []string{"Someone else <someone@fuad-daoud.com>", "Someone2 else2  <someone2@fuad-daoud.com>"},
		Pkgname:       "pkg-bin",
		Version:       "0.1.4",
		Pkgrel:        1,
		Description:   "Some single line description",
		Url:           "https://github.com/fuad-daoud/pkg",
		Arch:          []string{"x86_64", "aarch_64"},
		Licence:       []string{"MIT", "OBSD"},
		Provides:      []string{"package-a", "package-b"},
		Conflicts:     []string{"package-c", "package-d"},
		Source_x86_64: []string{"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64", "LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE", "README::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/README.md"},

		templatePath: "pkgbuild.tmpl",
	}

	result, err := data.template()

	if err != nil {
		fmt.Printf("got an err %v\n", err)
		t.FailNow()
	}

	expected, err := os.ReadFile("PKGBUILD_x86")

	if err != nil {
		fmt.Printf("got an err %v\n", err)
		t.FailNow()
	}
	assert.EqualValuesf(t, string(expected), result, "Failed Templating")
}
func TestValidate(t *testing.T) {
	for _, tt := range pkgbuild_tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pkg.validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
