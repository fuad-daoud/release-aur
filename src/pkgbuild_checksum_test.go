package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateForSourcesReal(t *testing.T) {
	sources := []string{"pkgmate-linux-amd64::https://github.com/fuad-daoud/pkgmate/releases/download/0.1.1/pkgmate-linux-amd64"}

	pkgbuild := &PkgBuild{
		client: NewClient(5*time.Second, 1*time.Second, 1),
	}

	result, err := pkgbuild.calculateForSources(sources)

	assert.NoError(t, err)
	// checksum from github
	assert.Equal(t, []string{"ba62160a8721ea41c112adc3cd369e1c7abb9d1c03d2bd89d13740b420cc1cc6"}, result)
}

func TestCalculateForSources(t *testing.T) {
	tests := []struct {
		name        string
		sources     []string
		serverSetup func() *httptest.Server
		wantErr     bool
		errContains string
		expected    []string
	}{

		{
			name:    "single source with known content",
			sources: []string{"https://example.com/file.tar.gz"},
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("test content"))
				}))
			},
			expected: []string{"6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"},
		},
		{
			name: "multiple sources",
			sources: []string{
				"file1::https://example.com/file1",
				"file2::https://example.com/file2",
			},
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/file1" {
						w.Write([]byte("content1"))
					} else if r.URL.Path == "/file2" {
						w.Write([]byte("content2"))
					}
				}))
			},
			expected: []string{
				"d0b425e00e15a0d36b9b361f02bab63563aed6cb4665083905386c55d5b679fa",
				"dab741b6289e7dccc1ed42330cae1accc2b755ce8079c2cd5d4b5366c9f769a6",
			},
		},
		{
			name:    "source with filename prefix",
			sources: []string{"myfile-1.0.0::https://example.com/release.tar.gz"},
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("release"))
				}))
			},
			expected: []string{"a4d451ec23463726f72c43d64c710968f6b602cd653b4de8adee1b556240a829"},
		},
		{
			name:    "download fails",
			sources: []string{"https://example.com/notfound"},
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			wantErr:     true,
			errContains: "failed to download",
		},
		{
			name:    "empty source list",
			sources: []string{},
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("shouldn't be called"))
				}))
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverSetup()
			defer server.Close()

			adjustedSources := make([]string, len(tt.sources))
			for i, source := range tt.sources {
				if idx := lastIndexOf(source, "::"); idx != -1 {
					prefix := source[:idx+2]
					path := source[idx+2:]
					if pathStart := lastIndexOf(path, "/"); pathStart != -1 {
						adjustedSources[i] = prefix + server.URL + path[pathStart:]
					} else {
						adjustedSources[i] = prefix + server.URL
					}
				} else {
					adjustedSources[i] = server.URL
				}
			}

			pkgbuild := &PkgBuild{
				client: Client{
					base:              server.URL,
					tries:             1,
					waitRetryDuration: 1 * time.Second,
					client: &http.Client{
						Timeout: 5 * time.Second,
					},
				},
			}

			result, err := pkgbuild.calculateForSources(adjustedSources)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func lastIndexOf(s, substr string) int {
	return strings.LastIndex(s, substr)
}
