package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCompareWithRemote_NewPackage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"resultcount":0,"results":[]}`))
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:         "new-package",
		Version:         "1.0.0",
		Checksum_x86_64: []string{"abc123"},
	}

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, "test PKGBUILD content")

	assert.NoError(t, err)
	assert.Equal(t, -1, pkgrel, "New package should return -1")
}

func TestDefaultCompareWithRemote_NewVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"0.9.0-1"}]}`))
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:         "test",
		Version:         "1.0.0",
		Checksum_x86_64: []string{"abc123"},
	}

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, "test PKGBUILD content")

	assert.NoError(t, err)
	assert.Equal(t, -1, pkgrel, "New version should return -1")
}

func TestDefaultCompareWithRemote_SameVersionDifferentContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-2"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			// Return different PKGBUILD
			w.Write([]byte(`pkgname=test
pkgver=1.0.0
pkgrel=2
sha256sums_x86_64=('abc123')`))
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:         "test",
		Version:         "1.0.0",
		Checksum_x86_64: []string{"abc123"},
	}

	localPKGBUILD := `pkgname=test
pkgver=1.0.0
pkgrel=1
description="Different description"
sha256sums_x86_64=('abc123')`

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.NoError(t, err)
	assert.Equal(t, 2, pkgrel, "Should increment pkgrel from remote")
}

func TestDefaultCompareWithRemote_SameVersionSameContent(t *testing.T) {
	localPKGBUILD := `pkgname=test
pkgver=1.0.0
pkgrel=1
sha256sums_x86_64=('abc123')`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			w.Write([]byte(localPKGBUILD))
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:         "test",
		Version:         "1.0.0",
		Checksum_x86_64: []string{"abc123"},
	}

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already published")
	assert.Equal(t, -1, pkgrel)
}

func TestDefaultCompareWithRemote_SameVersionX86_64NewChecksums(t *testing.T) {
	remotePKGBUILD := `pkgname=test
description="old description"
pkgver=1.0.0
sha256sums_x86_64=('oldchecksum')`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			w.Write([]byte(remotePKGBUILD))
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:         "test",
		Version:         "1.0.0",
		Checksum_x86_64: []string{"newchecksum"},
	}

	localPKGBUILD := `pkgname=test
description="new description"
pkgver=1.0.0
sha256sums_x86_64=('newchecksum')`

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "different x86_64 checksums")
	assert.Equal(t, -1, pkgrel)
}

func TestDefaultCompareWithRemote_SameVersionAARCH64NewChecksums(t *testing.T) {
	remotePKGBUILD := `pkgname=test
description="old description"
pkgver=1.0.0
sha256sums_aarch64=('oldchecksum')`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			w.Write([]byte(remotePKGBUILD))
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:          "test",
		Version:          "1.0.0",
		Checksum_aarch64: []string{"newchecksum"},
	}

	localPKGBUILD := `pkgname=test
description="new description"
pkgver=1.0.0
sha256sums_aarch64=('newchecksum')`

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "different aarch64 checksums")
	assert.Equal(t, -1, pkgrel)
}

func TestDefaultCompareWithRemote_MultipleArchitectures(t *testing.T) {
	remotePKGBUILD := `pkgname=test
pkgver=1.0.0
sha256sums_x86_64=('checksum1')
sha256sums_aarch64=('checksum2')`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-3"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			w.Write([]byte(remotePKGBUILD))
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:          "test",
		Version:          "1.0.0",
		Checksum_x86_64:  []string{"checksum1"},
		Checksum_aarch64: []string{"checksum2"},
	}

	localPKGBUILD := `pkgname=test
pkgver=1.0.0
description="Updated description"
sha256sums_x86_64=('checksum1')
sha256sums_aarch64=('checksum2')`

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.NoError(t, err)
	assert.Equal(t, 3, pkgrel)
}

func TestDefaultCompareWithRemote_FetchVersionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname: "test",
		Version: "1.0.0",
	}

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, "test")

	assert.Error(t, err)
	assert.Equal(t, -1, pkgrel)
}

func TestDefaultCompareWithRemote_FetchPKGBUILDError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:         "test",
		Version:         "1.0.0",
		Checksum_x86_64: []string{"abc"},
	}

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, "test")

	assert.Error(t, err)
	assert.Equal(t, -1, pkgrel)
}

func TestDefaultCompareWithRemote_Checksumx86_64CountMismatch(t *testing.T) {
	remotePKGBUILD := `pkgname=test
pkgver=1.0.0
sha256sums_x86_64=('checksum1' 'checksum2')`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			w.Write([]byte(remotePKGBUILD))
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:         "test",
		Version:         "1.0.0",
		Checksum_x86_64: []string{"checksum1"}, // Only one checksum
	}

	localPKGBUILD := `pkgname=test
pkgver=1.0.0
description="new description"
sha256sums_x86_64=('checksum1')`

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "different number of x86_64 checksums")
	assert.Equal(t, -1, pkgrel)
}

func TestDefaultCompareWithRemote_ChecksumAARCH64CountMismatch(t *testing.T) {
	remotePKGBUILD := `pkgname=test
pkgver=1.0.0
sha256sums_aarch64=('checksum1' 'checksum2')`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-1"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			w.Write([]byte(remotePKGBUILD))
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:          "test",
		Version:          "1.0.0",
		Checksum_aarch64: []string{"checksum1"}, // Only one checksum
	}

	localPKGBUILD := `pkgname=test
pkgver=1.0.0
description="new description"
sha256sums_aarch64=('checksum1')`

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "different number of aarch64 checksums")
	assert.Equal(t, -1, pkgrel)
}

func TestDefaultCompareWithRemote_RealPKGBUILD(t *testing.T) {
	aurPKGBUILD, err := os.ReadFile("testdata/PKGBUILD_AUR")
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-2"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			w.Write(aurPKGBUILD)
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:         "test",
		Version:         "1.0.0",
		Checksum_x86_64: []string{"SKIP"},
	}

	localPKGBUILD := string(aurPKGBUILD)

	_, err = defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already published")
}

func TestDefaultCompareWithRemote_SKIPChecksums(t *testing.T) {
	remotePKGBUILD := `pkgname=test
pkgver=1.0.0
sha256sums_x86_64=('SKIP')
sha256sums_aarch64=('SKIP')`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "rpc") {
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-3"}]}`))
		} else if strings.Contains(r.URL.Path, "PKGBUILD") {
			w.Write([]byte(remotePKGBUILD))
		}
	}))
	defer server.Close()

	client := DummyClient(server)
	pkgbuild := PkgBuild{
		Pkgname:          "test",
		Version:          "1.0.0",
		Checksum_x86_64:  []string{"checksum1"},
		Checksum_aarch64: []string{"checksum2"},
	}

	localPKGBUILD := `pkgname=test
pkgver=1.0.0
description="Updated description"
sha256sums_x86_64=('checksum1')
sha256sums_aarch64=('checksum2')`

	pkgrel, err := defaultCompareWithRemote(client, pkgbuild, localPKGBUILD)

	assert.NoError(t, err)
	assert.Equal(t, 3, pkgrel)
}
