package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAURClient_fetchPKGBUILD(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedContent := `pkgname=test
pkgver=1.0.0
pkgrel=1`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/cgit/aur.git/plain/PKGBUILD?h=test-pkg", r.URL.String())
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedContent))
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		result, err := client.fetchPKGBUILD("test-pkg")

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, result)
	})

	t.Run("404 not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		_, err := client.fetchPKGBUILD("nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("500 server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		_, err := client.fetchPKGBUILD("test-pkg")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})
}

func TestAURClient_getAurPackageVersions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/rpc/?v=5&type=info&arg[]=test-pkg", r.URL.String())
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test-pkg","Version":"1.2.3-5"}]}`))
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		result, err := client.getAurPackageVersions("test-pkg")

		assert.NoError(t, err)
		assert.Equal(t, "1.2.3", result.version)
		assert.Equal(t, 5, result.pkgrel)
	})

	t.Run("500 server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		_, err := client.getAurPackageVersions("test-pkg")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid json`))
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		_, err := client.getAurPackageVersions("test-pkg")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal")
	})

	t.Run("resultcount is 0", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"resultcount":0,"results":[]}`))
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		data, err := client.getAurPackageVersions("nonexistent")
		assert.NoError(t, err)
		assert.True(t, data.new)
	})

	t.Run("resultcount is more than 1", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"resultcount":2,"results":[{"Name":"test1","Version":"1.0.0-1"},{"Name":"test2","Version":"2.0.0-1"}]}`))
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		_, err := client.getAurPackageVersions("test")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid number of packages")
	})

	t.Run("non-numeric pkgrel", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"resultcount":1,"results":[{"Name":"test","Version":"1.0.0-abc"}]}`))
		}))
		defer server.Close()

		client := AURClient{base: server.URL}
		_, err := client.getAurPackageVersions("test")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse pkgRel")
	})
	t.Run("http.Get error in fetchPKGBUILD", func(t *testing.T) {
		client := AURClient{base: "http://invalid-url-that-does-not-exist", client: &http.Client{Timeout: 100 * time.Millisecond}}
		_, err := client.fetchPKGBUILD("test")

		assert.Error(t, err)
	})
	t.Run("http.Get error in getAurPackageVersions", func(t *testing.T) {
		client := AURClient{base: "http://invalid-url-that-does-not-exist", client: &http.Client{Timeout: 100 * time.Millisecond}}
		_, err := client.getAurPackageVersions("test")

		assert.Error(t, err)
	})
}
