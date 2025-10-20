package main

import (
	"fmt"
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

		pkgbuildTemplatePath: "pkgbuild.tmpl",
		outputPath:   "./output/PKGBUILD",
		client:       NewAURClient(5 * time.Second, 5 * time.Second, 5),
	}
	err := pkgbuild.validate()
	if err != nil {
		t.Error(fmt.Sprintf("got an err %v\n", err))
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
		pkgbuildTemplatePath: "pkgbuild.tmpl",
		client:       NewAURClient(5 * time.Second, 5 * time.Second, 5),
	}
	err := pkgbuild.validate()
	if err != nil {
		t.Error(fmt.Sprintf("got an err %v\n", err))
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
			pkgbuildTemplatePath:  "./nonexistent.tmpl",
			client:        NewAURClient(5*time.Second, time.Second, 1),
		}

		_, err := pkg.generate()
		assert.Error(t, err)
	})

	t.Run("getAurPackageVersions returns new", func(t *testing.T) {
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
			pkgbuildTemplatePath:  "./pkgbuild.tmpl",
			client:        DummyAURClient(server),
		}

		_, err := pkg.generate()
		assert.Error(t, err)
	})
	t.Run("getAurPackageVersions returns no package found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"resultcount":0,"results":[]}`))
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
			pkgbuildTemplatePath:  "./pkgbuild.tmpl",
			outputPath:    "/tmp/PKGBUILD",
			client:        DummyAURClient(server),
		}

		_, err := pkg.generate()
		assert.NoError(t, err)
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
			pkgbuildTemplatePath:  "./pkgbuild.tmpl",
			client:        DummyAURClient(server),
			outputPath:    "/root/PKGBUILD",
		}

		_, err := pkg.generate()
		assert.Error(t, err)
	})
	t.Run("PKGBUILDs match - already published", func(t *testing.T) {
		PKGBUILD_AUR, err := os.ReadFile("testdata/PKGBUILD_AUR")

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
			pkgbuildTemplatePath:  "./pkgbuild.tmpl",
			client:        DummyAURClient(server),
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
			pkgbuildTemplatePath:  "/tmp/pkgbuild.tmpl",
			client:        DummyAURClient(server),
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
			pkgbuildTemplatePath:  "pkgbuild.tmpl",
			client:        DummyAURClient(server),
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
