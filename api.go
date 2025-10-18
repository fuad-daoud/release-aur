package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

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
}

func fetchPKGBUILD(pkgName string) (string, error) {
	resp, err := http.Get("https://aur.archlinux.org/cgit/aur.git/plain/PKGBUILD?h=" + pkgName)

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

func getAurPackageVersions(pkgName string) (AurData, error) {
	resp, err := http.Get("https://aur.archlinux.org/rpc/?v=5&type=info&arg[]=" + pkgName)
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
	if result.Resultcount != 0 {
		return AurData{}, fmt.Errorf("Invalid number of packages in aur package: %s, found %v", pkgName, result.Results)
	}
	version := result.Results[0].Version
	index := strings.LastIndex(version, "-")
	pkgRel, err := strconv.Atoi(version[index:])
	if err != nil {
		return AurData{}, fmt.Errorf("Couldn't parse pkgRel in version %s", version[index:])
	}

	version = version[:index]
	return AurData{
		version: version,
		pkgrel:  pkgRel,
	}, nil
}
