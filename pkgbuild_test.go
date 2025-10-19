package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNewVersion(t *testing.T) {

	pkgbuild := PkgBuild{
		CliName:       "pkgmate",
		Maintainers:   []string{"Fuad Daoud <aur@fuad-daoud.com>"},
		Pkgname:       "pkgmate-bin",
		Version:       "100.0.0",
		Pkgrel:        1,
		Description:   "TUI application to manage your dependencies",
		Url:           "https://github.com/fuad-daoud/pkgmate",
		Arch:          []string{"x86_64"},
		Licence:       []string{"MIT"},
		Source_x86_64: []string{"pkgmate-bin-100.0.0-x86_64::https://github.com/fuad-daoud/pkgmate/releases/download/100.0.0/pkgmate-linux-amd64", "LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkgmate/v100.0.0/LICENSE", "README::https://raw.githubusercontent.com/fuad-daoud/pkgmate/v100.0.0/README.md"},

		templatePath: "pkgbuild.tmpl",
		outputPath:   "./output/PKGBUILD",
		client:       NewAURClient(5 * time.Second),
	}
	err := pkgbuild.validate()
	if err != nil {
		t.Error(fmt.Sprintf("got an err %v\n", err))
	}

	PKGBUILD, err := pkgbuild.generate()
	if err != nil {
		t.Errorf("got an err %v", err)
	}

	expected, err := os.ReadFile("expected/PKGBUILD_pkgmate")

	if err != nil {
		t.Errorf("got an err %v\n", err)
	}
	assert.EqualValuesf(t, string(expected), PKGBUILD, "Failed Templating")

}

func TestGenerateNewPkgrel(t *testing.T) {
	pkgbuild := PkgBuild{
		CliName:       "pkgmate",
		Maintainers:   []string{"Fuad Daoud <aur@fuad-daoud.com>"},
		Pkgname:       "pkgmate-bin",
		Version:       "0.1.1",
		Pkgrel:        1,
		Description:   "TUI application to manage your dependencies",
		Url:           "https://github.com/fuad-daoud/pkgmate",
		Arch:          []string{"x86_64"},
		Licence:       []string{"MIT"},
		Source_x86_64: []string{"pkgmate-bin-0.1.1-x86_64::https://github.com/fuad-daoud/pkgmate/releases/download/0.1.1/pkgmate-linux-amd64", "LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkgmate/v0.1.1/LICENSE", "README::https://raw.githubusercontent.com/fuad-daoud/pkgmate/v0.1.1/README.md"},

		outputPath:   "./output/PKGBUILD",
		templatePath: "pkgbuild.tmpl",
		client:       NewAURClient(5 * time.Second),
	}
	err := pkgbuild.validate()
	if err != nil {
		t.Error(fmt.Sprintf("got an err %v\n", err))
	}

	PKGBUILD, err := pkgbuild.generate()
	if err != nil {
		t.Errorf("got an err %v", err)
	}

	expected, err := os.ReadFile("expected/PKGBUILD_pkgmate_rel2")

	if err != nil {
		t.Errorf("got an err %v\n", err)
	}
	assert.EqualValuesf(t, string(expected), PKGBUILD, "Failed Templating")

}
func TestGenerate_Errors_WithHttpTest(t *testing.T) {

	t.Run("template fails", func(t *testing.T) {
		pkg := &PkgBuild{
			CliName:       "test",
			Maintainers:   []string{"User"},
			Pkgname:       "test",
			Version:       "1.0.0",
			Description:   "Test",
			Url:           "https://example.com",
			Arch:          []string{"x86_64"},
			Licence:       []string{"MIT"},
			Source_x86_64: []string{"test"},
			templatePath:  "./nonexistent.tmpl",
			client:        NewAURClient(5 * time.Second),
		}

		_, err := pkg.generate()
		assert.Error(t, err)
	})
	t.Run("getAurPackageVersions returns 500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		pkg := &PkgBuild{
			CliName:       "test",
			Maintainers:   []string{"User"},
			Pkgname:       "test",
			Version:       "1.0.0",
			Description:   "Test",
			Url:           "https://example.com",
			Arch:          []string{"x86_64"},
			Licence:       []string{"MIT"},
			Source_x86_64: []string{"test"},
			templatePath:  "./pkgbuild.tmpl",
			client:        AURClient{base: server.URL},
		}

		_, err := pkg.generate()
		assert.Error(t, err)
	})

	t.Run("fetchPKGBUILD fails when versions match", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "rpc") {
				w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		pkg := &PkgBuild{
			CliName:       "test",
			Maintainers:   []string{"User"},
			Pkgname:       "test",
			Version:       "1.0.0",
			Description:   "Test",
			Url:           "https://example.com",
			Arch:          []string{"x86_64"},
			Licence:       []string{"MIT"},
			Source_x86_64: []string{"test"},
			templatePath:  "./pkgbuild.tmpl",
			client:        AURClient{base: server.URL},
			outputPath:    "/root/PKGBUILD",
		}

		_, err := pkg.generate()
		assert.Error(t, err)
	})
	t.Run("PKGBUILDs match - already published", func(t *testing.T) {

		PKGBUILD_AUR, err := os.ReadFile("expected/PKGBUILD_AUR")

		if err != nil {
			t.Errorf("got an err %v\n", err)
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "rpc") {
				w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
			} else {
				w.Write([]byte(PKGBUILD_AUR))

			}
		}))
		defer server.Close()

		pkg := &PkgBuild{
			CliName:       "test",
			Maintainers:   []string{"User"},
			Pkgname:       "test",
			Version:       "1.0.0",
			Description:   "Test",
			Url:           "https://example.com",
			Arch:          []string{"x86_64"},
			Licence:       []string{"MIT"},
			Source_x86_64: []string{"test"},
			templatePath:  "./pkgbuild.tmpl",
			client:        AURClient{base: server.URL},
			outputPath:    "/root/PKGBUILD",
		}

		PKGBUILD, err := pkg.generate()
		if err == nil {
			fmt.Printf("PKGBUILD: \n%v\n", PKGBUILD)

			assert.EqualValuesf(t, PKGBUILD_AUR, PKGBUILD, "Failed Templating")
		}
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already published")
	})

	t.Run("PKGBUILDs match - new pkgrel, second templating fails", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "rpc") {
				w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
			} else {
				w.Write([]byte("different content"))

				err := copyFile("invalid_template.tmpl", "/tmp/pkgbuild.tmpl")
				if err != nil {
					t.Errorf("error in setup, %v", err)
				}
			}
		}))
		err := copyFile("pkgbuild.tmpl", "/tmp/pkgbuild.tmpl")
		if err != nil {
			t.Errorf("error in setup, %v", err)
		}

		pkg := &PkgBuild{
			CliName:       "test",
			Maintainers:   []string{"User"},
			Pkgname:       "test",
			Version:       "1.0.0",
			Description:   "Test",
			Url:           "https://example.com",
			Arch:          []string{"x86_64"},
			Licence:       []string{"MIT"},
			Source_x86_64: []string{"test"},
			templatePath:  "/tmp/pkgbuild.tmpl",
			client:        AURClient{base: server.URL},
			outputPath:    "/root/PKGBUILD",
		}

		_, err = pkg.generate()
		if err == nil {
			t.Errorf("should have retruned an error")
		}
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template")
	})

	t.Run("Geneartion works, fails at writing", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "rpc") {
				w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.1-1"}]}`))
			}
		}))
		pkg := &PkgBuild{
			CliName:       "test",
			Maintainers:   []string{"User"},
			Pkgname:       "test",
			Version:       "1.0.0",
			Description:   "Test",
			Url:           "https://example.com",
			Arch:          []string{"x86_64"},
			Licence:       []string{"MIT"},
			Source_x86_64: []string{"test"},
			templatePath:  "pkgbuild.tmpl",
			client:        AURClient{base: server.URL},
			outputPath:    "/root/PKGBUILD",
		}

		_, err := pkg.generate()
		if err == nil {
			t.Errorf("should have retruned an error")
		}
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")
	})
}

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
		t.Errorf("got an err %v\n", err)
	}

	expected, err := os.ReadFile("expected/PKGBUILD")

	if err != nil {
		t.Errorf("got an err %v\n", err)
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
		t.Errorf("got an err %v\n", err)
	}

	expected, err := os.ReadFile("expected/PKGBUILD_x86")

	if err != nil {
		t.Errorf("got an err %v\n", err)
	}
	assert.EqualValuesf(t, string(expected), result, "Failed Templating")
}
func TestValidate(t *testing.T) {
	for _, tt := range pkgbuildTests {
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
	for _, tt := range pkgbuildFromEnvTests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			result := NewPkgBuildFromEnv()

			assert.Equal(t, tt.expected.Maintainers, result.Maintainers)
			assert.Equal(t, tt.expected.Contributors, result.Contributors)
			assert.Equal(t, tt.expected.Pkgname, result.Pkgname)
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

func TestTemplate_Errors(t *testing.T) {
	tests := []struct {
		name    string
		pkg     PkgBuild
		wantErr bool
	}{
		{
			name: "invalid template path",
			pkg: PkgBuild{
				templatePath: "./nonexistent.tmpl",
			},
			wantErr: true,
		},
		{
			name: "empty template path",
			pkg: PkgBuild{
				templatePath: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid template",
			pkg: PkgBuild{
				templatePath: "./invalid_template.tmpl",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.pkg.template()
			assert.Error(t, err)
		})
	}
}

func TestWriteFile_Errors(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		content  string
		setup    func()
		cleanup  func()
	}{
		{
			name:     "invalid path - read-only parent",
			filePath: "/root/test/PKGBUILD",
			content:  "test",
		},
		{
			name:     "permission denied on directory creation",
			filePath: "/test-readonly/subdir/file.txt",
			content:  "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			err := writeFile(tt.filePath, tt.content)
			assert.Error(t, err)
		})
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
