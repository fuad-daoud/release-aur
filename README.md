# Archiving in favor of using gorelaser


# Release to AUR GitHub Action

[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Release%20to%20AUR-blue.svg?colorA=24292e&colorB=0366d6&style=flat&longCache=true&logo=github)](https://github.com/marketplace/actions/release-to-aur)
[![CI](https://github.com/fuad-daoud/release-aur/workflows/Test%20Action/badge.svg)](https://github.com/fuad-daoud/release-aur/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A GitHub Action that generates PKGBUILD files for publishing packages to the Arch User Repository (AUR).

## Features

- 🚀 Automatically generates PKGBUILD files from release artifacts
- 🔄 Handles version management and pkgrel increments
- 🏗️ Supports multiple architectures (x86_64, aarch64)
- ✅ Validates against existing AUR packages
- 📦 Compares with published PKGBUILDs to avoid duplicates

## Usage

### Basic Example

```yaml
name: Release to AUR

on:
  release:
    types: [published]

jobs:
  aur-release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Generate PKGBUILD
        uses: fuad-daoud/release-aur@latest
        with:
          cli_name: 'myapp'
          maintainers: 'Your Name <your.email@example.com>'
          pkgname: 'myapp-bin'
          version: ${{ github.event.release.tag_name }}
          description: 'My awesome application'
          url: 'https://github.com/${{ github.repository }}'
          arch: 'x86_64,aarch64'
          licence: 'MIT'
          source_x86_64: 'https://github.com/${{ github.repository }}/releases/download/${{ github.event.release.tag_name }}/myapp-linux-amd64'
          source_aarch64: 'https://github.com/${{ github.repository }}/releases/download/${{ github.event.release.tag_name }}/myapp-linux-arm64'

      - name: Upload PKGBUILD
        uses: actions/upload-artifact@v4
        with:
          name: pkgbuild
          path: output/PKGBUILD
```

### Advanced Example with AUR Publishing

```yaml
name: Release to AUR

on:
  release:
    types: [published]

jobs:
  aur-release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Generate PKGBUILD
        uses: fuad-daoud/release-aur@v1
        with:
          cli_name: 'myapp'
          maintainers: 'Your Name <your.email@example.com>'
          contributors: 'Contributor Name <contrib@example.com>'
          pkgname: 'myapp-bin'
          version: ${{ github.event.release.tag_name }}
          description: 'My awesome application'
          url: 'https://github.com/${{ github.repository }}'
          arch: 'x86_64,aarch64'
          licence: 'MIT'
          provides: 'myapp'
          conflicts: 'myapp-git'
          source_x86_64: 'https://github.com/${{ github.repository }}/releases/download/${{ github.event.release.tag_name }}/myapp-linux-amd64'
          source_aarch64: 'https://github.com/${{ github.repository }}/releases/download/${{ github.event.release.tag_name }}/myapp-linux-arm64'

      - name: Publish to AUR
        uses: KSXGitHub/github-actions-deploy-aur@v2
        with:
          pkgname: myapp-bin
          pkgbuild: output/PKGBUILD
          commit_username: ${{ secrets.AUR_USERNAME }}
          commit_email: ${{ secrets.AUR_EMAIL }}
          ssh_private_key: ${{ secrets.AUR_SSH_PRIVATE_KEY }}
          commit_message: "Update to version ${{ github.event.release.tag_name }}"
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `cli_name` | Name of the CLI binary to install | Yes | - |
| `maintainers` | Comma-separated list of maintainers | Yes | - |
| `contributors` | Comma-separated list of contributors | No | `''` |
| `pkgname` | Package name for AUR | Yes | - |
| `version` | Version of the package | Yes | - |
| `description` | Package description | Yes | - |
| `url` | Project URL | Yes | - |
| `arch` | Comma-separated list of architectures | Yes | - |
| `licence` | Comma-separated list of licenses | Yes | - |
| `provides` | Comma-separated list of provided packages | No | `''` |
| `conflicts` | Comma-separated list of conflicting packages | No | `''` |
| `source_x86_64` | Comma-separated list of x86_64 source URLs | Yes | - |
| `source_aarch64` | Comma-separated list of aarch64 source URLs | No | `''` |
| `pkgbuild_template` | Path to custom PKGBUILD template relative to the github action path | No | `src/pkgbuild.tmpl` |
| `srcinfo_template` | Path to custom .SRCINFO template relative to the github action path | No | `src/srcinfo.tmpl` |
| `output_path` | Output path where the PKGBUILD will be generated relative to workspace root | No | `PKGBUILD` |

## Outputs

| Output | Description |
|--------|-------------|
| `pkgbuild_path` | Path to the generated PKGBUILD file |

## How It Works

1. **Validation**: Validates all required inputs
2. **AUR Check**: Fetches current version from AUR (if exists)
3. **Version Comparison**: 
   - If version matches: Compares PKGBUILD content and increments `pkgrel`
   - If version differs: Resets `pkgrel` to 1
4. **Generation**: Creates PKGBUILD from template
5. **Output**: Saves PKGBUILD to specified path

## Development

### Running Tests

```bash
go test -v ./...
```

### Running Locally

```bash
export cli_name="myapp"
export maintainers="Your Name <email@example.com>"
export pkgname="myapp-bin"
export version="1.0.0"
export description="My app description"
export url="https://github.com/user/repo"
export arch="x86_64"
export licence="MIT"
export source_x86_64="https://example.com/myapp-1.0.0-x86_64"

go run .
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
