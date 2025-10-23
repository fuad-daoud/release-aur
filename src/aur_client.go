package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	base              string
	client            *http.Client
	tries             int
	waitRetryDuration time.Duration
}

func NewClient(timeout, waitRetryDuration time.Duration, tries int) Client {
	return Client{
		base:              "https://aur.archlinux.org",
		tries:             max(tries, 1),
		waitRetryDuration: max(100*time.Millisecond, waitRetryDuration),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type AurResponse struct {
	Resultcount int `json:"resultcount"`
	Results     []struct {
		Name    string `json:"Name"`
		Version string `json:"Version"`
	} `json:"results"`
}
type AurData struct {
	version string
	pkgrel  int
	new     bool
}

func (client Client) Get(url string) ([]byte, error) {
	if client.client == nil {
		client.client = &http.Client{Timeout: 30 * time.Second}
	}

	var resp *http.Response
	var err error
	for client.tries > 0 {
		resp, err = client.client.Get(url)

		if err != nil {
			return []byte{}, err
		}

		if resp.StatusCode == 200 {
			break
		}
		client.tries--
		slog.Warn("Got wrong status trying again", "duration before retry", client.waitRetryDuration, "tries left", client.tries)
		if client.tries == 0 {
			break
		}
		time.Sleep(client.waitRetryDuration)
	}
	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("Error integrating got none 200 status %v\n", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (client Client) getAur(path string) ([]byte, error) {
	return client.Get(client.base + path)
}

func (client Client) fetchPKGBUILD(pkgName string) (string, error) {
	body, err := client.getAur("/cgit/aur.git/plain/PKGBUILD?h=" + pkgName)
	return string(body), err
}

func (client Client) getAurPackageVersions(pkgName string) (AurData, error) {
	body, err := client.getAur("/rpc/?v=5&type=info&arg[]=" + pkgName)
	if err != nil {
		return AurData{}, err
	}
	var result AurResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		slog.Error("could not parse", "jsonString", string(body))
		return AurData{}, fmt.Errorf("Could not unmarshal the response: %v\n", err)
	}
	if result.Resultcount > 1 {
		return AurData{}, fmt.Errorf("Invalid number of packages in aur package: %s, found %v", pkgName, result.Results)
	}
	if result.Resultcount == 0 {
		return AurData{new: true}, nil
	}
	version := result.Results[0].Version
	index := strings.LastIndex(version, "-")
	pkgRel, err := strconv.Atoi(version[index+1:])
	if err != nil {
		return AurData{}, fmt.Errorf("Couldn't parse pkgRel in version %s", version[index:])
	}

	version = version[:index]
	return AurData{
		version: version,
		pkgrel:  pkgRel,
	}, nil
}
