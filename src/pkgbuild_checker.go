package main

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/fuad-daoud/release-aur/src/parser"
)

func validate(p PkgBuild) error {
	slog.Info("Validating", "pkgbuild", p)
	if p.CliName == "" {
		return fmt.Errorf("CliName is required")
	}
	if len(p.Maintainers) == 0 {
		return fmt.Errorf("At least one Maintainer is required")
	}
	if p.Pkgname == "" {
		return fmt.Errorf("Pkgname is required")
	}
	if p.Version == "" {
		return fmt.Errorf("Version is required")
	}
	if p.Description == "" {
		return fmt.Errorf("Description is required")
	}
	if p.Url == "" {
		return fmt.Errorf("Url is required")
	}
	if len(p.Arch) == 0 {
		return fmt.Errorf("At least one Arch is required")
	}
	if len(p.Licence) == 0 {
		return fmt.Errorf("At least one Licence is required")
	}
	if len(p.Source_x86_64) == 0 {
		return fmt.Errorf("Source_x86_64 is required")
	}
	slog.Info("Input is valid")
	return nil
}

type compareWithRemote func(client Client, pkgbuild PkgBuild, PKGBUILD string) (int, error)

func defaultCompareWithRemote(client Client, pkgbuild PkgBuild, PKGBUILD string) (int, error) {

	data, err := client.getAurPackageVersions(pkgbuild.Pkgname)

	if err != nil {
		slog.Error("Failed to fetch package info from AUR")
		return -1, err
	}

	if data.new == false && data.version == pkgbuild.Version {
		slog.Warn("AUR version and current version match, this should only be a PKGBUILD update")
		slog.Info("Comparing PKGBUILD to validate")
		aurPKGBUILD, err := client.fetchPKGBUILD(pkgbuild.Pkgname)
		if err != nil {
			slog.Error("Failed to fetch PKGBUILD from AUR")
			return -1, err
		}
		remoteChecksums, err := parser.ExtractChecksums(PKGBUILD)

		if err != nil {
			slog.Error("Failed to fetch PKGBUILD from AUR")
			return -1, err
		}
		if parser.ComparePKGBUILDs(PKGBUILD, aurPKGBUILD) {
			slog.Error("Files match!! should not publish to the AUR without changes to PKGBUILD file or the software Version")
			return -1, fmt.Errorf("PKGBUILD already published to AUR")
		}

		if len(pkgbuild.Checksum_aarch64)+len(pkgbuild.Checksum_x86_64) != len(remoteChecksums) {
			slog.Error("Checksums differ!! should not increament the pkgrel with new sources")
			return -1, fmt.Errorf("Old version new checksums")
		}
		for _, local := range pkgbuild.Checksum_aarch64 {
			found := slices.Contains(remoteChecksums["aarch64"], local)
			if !found {
				slog.Error("Checksums differ!! should not increament the pkgrel with new sources")
				return -1, fmt.Errorf("Old version new checksums")
			}
		}

		for _, local := range pkgbuild.Checksum_x86_64 {
			found := slices.Contains(remoteChecksums["x86_64"], local)
			if !found {
				slog.Error("Checksums differ!! should not increament the pkgrel with new sources")
				return -1, fmt.Errorf("Old version new checksums")
			}
		}
		return data.pkgrel, nil
	}
	if data.new {
		slog.Info("New package")
	}
	return -1, nil
}
