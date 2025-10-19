package main

var normalize_tests = []struct {
	name     string
	input    string
	expected string
}{
	{
		name: "removes checksums and pkgrel",
		input: `pkgname=test
pkgver=1.0.0
pkgrel=1
sha256sums=('abc123')
sha256sums_x86_64=('def456')`,
		expected: `pkgname=test
pkgver=1.0.0`,
	},
	{
		name: "removes multi-line checksum arrays",
		input: `pkgname=test
pkgver=1.0.0
pkgrel=2
sha256sums_x86_64=('7f6d936dae7da64b45acdadef654863f2d47866660c6cf821430707bbf63c4cd'
                   'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4'
                   '8320695d303094310734f5df7e96722a03d7948b076d52849a7b014006aff793')
sha256sums_aarch64=('3bc82df11d552c8134436fd0d752cda0f551932d9e80f038110f0c3ec0e39232'
                    'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4'
                    '8320695d303094310734f5df7e96722a03d7948b076d52849a7b014006aff793')
arch=('x86_64')`,
		expected: `pkgname=test
pkgver=1.0.0
arch=('x86_64')`,
	},
	{
		name: "keeps source arrays",
		input: `pkgname=test
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
pkgrel=1
sha256sums=('SKIP')`,
		expected: `pkgname=test
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)`,
	},
	{
		name: "handles SKIP checksums",
		input: `pkgname=test
sha256sums=('SKIP')
sha256sums_x86_64=('SKIP')`,
		expected: `pkgname=test`,
	},
	{
		name: "full PKGBUILD example",
		input: `# Maintainer: Test User
pkgname=prayers-bin
pkgver=0.1.4
pkgrel=1
pkgdesc="Test package"
arch=('x86_64' 'aarch64')
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
source_aarch64=(
"pkg-bin-0.1.4-aarch_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-aarch_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
sha256sums_x86_64=('7f6d936dae7da64b45acdadef654863f2d47866660c6cf821430707bbf63c4cd'
                   'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4')
sha256sums_aarch64=('3bc82df11d552c8134436fd0d752cda0f551932d9e80f038110f0c3ec0e39232'
                    'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4')
package() {
    install -Dm755 "$srcdir/pkg-bin" "$pkgdir/usr/bin/prayers"
}`,
		expected: `# Maintainer: Test User
pkgname=prayers-bin
pkgver=0.1.4
pkgdesc="Test package"
arch=('x86_64' 'aarch64')
source_x86_64=(
"pkg-bin-0.1.4-x86_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-x86_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
source_aarch64=(
"pkg-bin-0.1.4-aarch_64::https://github.com/fuad-daoud/pkg/releases/download/v0.1.4/prayers-linux-aarch_64"
"LICENSE::https://raw.githubusercontent.com/fuad-daoud/pkg/v0.1.4/LICENSE"
)
package() {
    install -Dm755 "$srcdir/pkg-bin" "$pkgdir/usr/bin/prayers"
}`,
	},
}
