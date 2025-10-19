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

type AURClient struct {
	base   string
	client *http.Client
}

func NewAURClient(timeout time.Duration) AURClient {
	return AURClient{
		base: "https://aur.archlinux.org",
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

func (client AURClient) fetchPKGBUILD(pkgName string) (string, error) {
	if client.client == nil {
		client.client = &http.Client{Timeout: 10 * time.Second}
	}
	resp, err := client.client.Get(client.base + "/cgit/aur.git/plain/PKGBUILD?h=" + pkgName)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error integrating got none 200 status %v\n", resp.StatusCode)
	}

	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func (client AURClient) getAurPackageVersions(pkgName string) (AurData, error) {
	if client.client == nil {
		client.client = &http.Client{Timeout: 10 * time.Second}
	}
	resp, err := client.client.Get(client.base + "/rpc/?v=5&type=info&arg[]=" + pkgName)
	if err != nil {
		return AurData{}, err
	}

	if resp.StatusCode != 200 {
		return AurData{}, fmt.Errorf("Error integrating got none 200 status %v\n", resp.StatusCode)
	}
	defer resp.Body.Close()

	jsonString, err := io.ReadAll(resp.Body)
	if err != nil {
		return AurData{}, err
	}
	var result AurResponse
	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		slog.Error("could not parse", "jsonString", jsonString)
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
