#!/bin/zsh

set -e

/bin/rm -rf payload
/bin/rm -rf build
/bin/rm -rf output
/bin/mkdir -p payload/opt/macadmins/bin
/bin/mkdir -p build
/bin/mkdir -p output

APP_SIGNING_IDENTITY="Developer ID Application: Mac Admins Open Source (T4SK8ZXCXG)"
INSTALLER_SIGNING_IDENTITY="Developer ID Installer: Mac Admins Open Source (T4SK8ZXCXG)"

# read the version from the VERSION file
VERSION=$(cat VERSION)

# build arm64
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o build/default-browser-arm64

# build amd64
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o build/default-browser-amd64

# combine the two into a universal binary
/usr/bin/lipo -create -output build/default-browser build/default-browser-arm64 build/default-browser-amd64

# copy the binary to the payload
/bin/cp build/default-browser payload/opt/macadmins/bin/default-browser

# sign the binary
sudo /usr/bin/codesign --timestamp --force --deep -s "${APP_SIGNING_IDENTITY}" payload/opt/macadmins/bin/default-browser

# create the package
/usr/bin/pkgbuild --root payload --identifier com.github.macadmins.default-browser --version ${VERSION} --install-location / --ownership recommended --sign "${INSTALLER_SIGNING_IDENTITY}" output/default-browser.pkg