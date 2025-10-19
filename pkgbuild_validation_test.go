package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		pkg     PkgBuild
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid package",
			pkg: PkgBuild{
				CliName:       "test",
				Maintainers:   []string{"Test User <test@example.com>"},
				Pkgname:       "test-bin",
				Version:       "1.0.0",
				Description:   "Test package",
				Url:           "https://example.com",
				Arch:          []string{"x86_64"},
				Licence:       []string{"MIT"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: false,
		},
		{
			name: "missing CliName",
			pkg: PkgBuild{
				Maintainers:   []string{"Test User"},
				Pkgname:       "test-bin",
				Version:       "1.0.0",
				Description:   "Test package",
				Url:           "https://example.com",
				Arch:          []string{"x86_64"},
				Licence:       []string{"MIT"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: true,
			errMsg:  "CliName is required",
		},
		{
			name: "missing Maintainers",
			pkg: PkgBuild{
				CliName:       "test",
				Pkgname:       "test-bin",
				Version:       "1.0.0",
				Description:   "Test package",
				Url:           "https://example.com",
				Arch:          []string{"x86_64"},
				Licence:       []string{"MIT"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: true,
			errMsg:  "At least one Maintainer is required",
		},
		{
			name: "missing Pkgname",
			pkg: PkgBuild{
				CliName:       "test",
				Maintainers:   []string{"Test User"},
				Version:       "1.0.0",
				Description:   "Test package",
				Url:           "https://example.com",
				Arch:          []string{"x86_64"},
				Licence:       []string{"MIT"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: true,
			errMsg:  "Pkgname is required",
		},
		{
			name: "missing Version",
			pkg: PkgBuild{
				CliName:       "test",
				Maintainers:   []string{"Test User"},
				Pkgname:       "test-bin",
				Description:   "Test package",
				Url:           "https://example.com",
				Arch:          []string{"x86_64"},
				Licence:       []string{"MIT"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: true,
			errMsg:  "Version is required",
		},
		{
			name: "missing Description",
			pkg: PkgBuild{
				CliName:       "test",
				Maintainers:   []string{"Test User"},
				Pkgname:       "test-bin",
				Version:       "1.0.0",
				Url:           "https://example.com",
				Arch:          []string{"x86_64"},
				Licence:       []string{"MIT"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: true,
			errMsg:  "Description is required",
		},
		{
			name: "missing Url",
			pkg: PkgBuild{
				CliName:       "test",
				Maintainers:   []string{"Test User"},
				Pkgname:       "test-bin",
				Version:       "1.0.0",
				Description:   "Test package",
				Arch:          []string{"x86_64"},
				Licence:       []string{"MIT"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: true,
			errMsg:  "Url is required",
		},
		{
			name: "missing Arch",
			pkg: PkgBuild{
				CliName:       "test",
				Maintainers:   []string{"Test User"},
				Pkgname:       "test-bin",
				Version:       "1.0.0",
				Description:   "Test package",
				Url:           "https://example.com",
				Licence:       []string{"MIT"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: true,
			errMsg:  "At least one Arch is required",
		},
		{
			name: "missing Licence",
			pkg: PkgBuild{
				CliName:       "test",
				Maintainers:   []string{"Test User"},
				Pkgname:       "test-bin",
				Version:       "1.0.0",
				Description:   "Test package",
				Url:           "https://example.com",
				Arch:          []string{"x86_64"},
				Source_x86_64: []string{"https://example.com/test"},
			},
			wantErr: true,
			errMsg:  "At least one Licence is required",
		},
		{
			name: "missing Source_x86_64",
			pkg: PkgBuild{
				CliName:     "test",
				Maintainers: []string{"Test User"},
				Pkgname:     "test-bin",
				Version:     "1.0.0",
				Description: "Test package",
				Url:         "https://example.com",
				Arch:        []string{"x86_64"},
				Licence:     []string{"MIT"},
			},
			wantErr: true,
			errMsg:  "Source_x86_64 is required",
		},
		{
			name: "optional fields can be empty",
			pkg: PkgBuild{
				CliName:        "test",
				Maintainers:    []string{"Test User"},
				Pkgname:        "test-bin",
				Version:        "1.0.0",
				Description:    "Test package",
				Url:            "https://example.com",
				Arch:           []string{"x86_64"},
				Licence:        []string{"MIT"},
				Source_x86_64:  []string{"https://example.com/test"},
				Contributors:   nil,
				Provides:       nil,
				Conflicts:      nil,
				Source_aarch64: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
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
func TestNewPkgBuildFromEnv(t *testing.T) {

	tests := []struct {
		name     string
		envVars  map[string]string
		expected PkgBuild
	}{
		{
			name: "all fields provided",
			envVars: map[string]string{
				"maintainers":       "User1 <user1@example.com>,User2 <user2@example.com>",
				"contributors":      "Contrib1 <c1@example.com>,Contrib2 <c2@example.com>",
				"pkgname":           "test-bin",
				"cli_name":          "test",
				"version":           "1.0.0",
				"description":       "Test package",
				"url":               "https://example.com",
				"arch":              "x86_64,aarch64",
				"licence":           "MIT,Apache",
				"provides":          "test,test-cli",
				"conflicts":         "old-test",
				"source_x86_64":     "https://example.com/x86",
				"source_aarch64":    "https://example.com/arm",
				"pkgbuild_template": "./custom.tmpl",
			},
			expected: PkgBuild{
				Maintainers:    []string{"User1 <user1@example.com>", "User2 <user2@example.com>"},
				Contributors:   []string{"Contrib1 <c1@example.com>", "Contrib2 <c2@example.com>"},
				CliName:        "test",
				Pkgname:        "test-bin",
				Version:        "1.0.0",
				Pkgrel:         1,
				Description:    "Test package",
				Url:            "https://example.com",
				Arch:           []string{"x86_64", "aarch64"},
				Licence:        []string{"MIT", "Apache"},
				Provides:       []string{"test", "test-cli"},
				Conflicts:      []string{"old-test"},
				Source_x86_64:  []string{"https://example.com/x86"},
				Source_aarch64: []string{"https://example.com/arm"},
				templatePath:   "./custom.tmpl",
			},
		},
		{
			name: "optional fields empty",
			envVars: map[string]string{
				"maintainers":   "User1 <user1@example.com>",
				"pkgname":       "test-bin",
				"version":       "1.0.0",
				"description":   "Test package",
				"url":           "https://example.com",
				"arch":          "x86_64",
				"licence":       "MIT",
				"source_x86_64": "https://example.com/x86",
			},
			expected: PkgBuild{
				Maintainers:    []string{"User1 <user1@example.com>"},
				Contributors:   []string{},
				Pkgname:        "test-bin",
				Version:        "1.0.0",
				Pkgrel:         1,
				Description:    "Test package",
				Url:            "https://example.com",
				Arch:           []string{"x86_64"},
				Licence:        []string{"MIT"},
				Provides:       []string{},
				Conflicts:      []string{},
				Source_x86_64:  []string{"https://example.com/x86"},
				Source_aarch64: []string{},
				templatePath:   "./pkgbuild.tmpl",
			},
		},
		{
			name: "default template path when not provided",
			envVars: map[string]string{
				"maintainers":   "User1",
				"pkgname":       "test",
				"version":       "1.0.0",
				"description":   "Test",
				"url":           "https://example.com",
				"arch":          "x86_64",
				"licence":       "MIT",
				"source_x86_64": "https://example.com/x86",
			},
			expected: PkgBuild{
				Maintainers:    []string{"User1"},
				Contributors:   []string{},
				Pkgname:        "test",
				Version:        "1.0.0",
				Pkgrel:         1,
				Description:    "Test",
				Url:            "https://example.com",
				Arch:           []string{"x86_64"},
				Licence:        []string{"MIT"},
				Provides:       []string{},
				Conflicts:      []string{},
				Source_x86_64:  []string{"https://example.com/x86"},
				Source_aarch64: []string{},
				templatePath:   "./pkgbuild.tmpl",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			result := NewPkgBuildFromEnv()

			assert.Equal(t, tt.expected.Maintainers, result.Maintainers)
			assert.Equal(t, tt.expected.Contributors, result.Contributors)
			assert.Equal(t, tt.expected.Pkgname, result.Pkgname)
			assert.Equal(t, tt.expected.CliName, result.CliName)
			assert.Equal(t, tt.expected.Version, result.Version)
			assert.Equal(t, tt.expected.Pkgrel, result.Pkgrel)
			assert.Equal(t, tt.expected.Description, result.Description)
			assert.Equal(t, tt.expected.Url, result.Url)
			assert.Equal(t, tt.expected.Arch, result.Arch)
			assert.Equal(t, tt.expected.Licence, result.Licence)
			assert.Equal(t, tt.expected.Provides, result.Provides)
			assert.Equal(t, tt.expected.Conflicts, result.Conflicts)
			assert.Equal(t, tt.expected.Source_x86_64, result.Source_x86_64)
			assert.Equal(t, tt.expected.Source_aarch64, result.Source_aarch64)
			assert.Equal(t, tt.expected.templatePath, result.templatePath)
		})
	}
}
