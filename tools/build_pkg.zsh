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
XCODE_NOTARY_PATH="$XCODE_PATH/Contents/Developer/usr/bin/notarytool"
XCODE_STAPLER_PATH="$XCODE_PATH/Contents/Developer/usr/bin/stapler"

# read the version from passed argument
VERSION=$1

echo "Building version ${VERSION}"

# build arm64
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o build/default-browser-arm64

# build amd64
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o build/default-browser-amd64

# combine the two into a universal binary
echo "Combining arm64 and amd64 into a universal binary"
/usr/bin/lipo -create -output build/default-browser build/default-browser-arm64 build/default-browser-amd64

# copy the binary to the payload
echo "Copying the binary to the payload"
/bin/cp build/default-browser payload/opt/macadmins/bin/default-browser

# sign the binary
echo "Signing the binary"
sudo /usr/bin/codesign --timestamp --force --deep -s "${APP_SIGNING_IDENTITY}" payload/opt/macadmins/bin/default-browser

# create the package
echo "Creating the package"
/usr/bin/pkgbuild --root payload --identifier com.github.macadmins.default-browser --version ${VERSION} --install-location / --ownership recommended --sign "${INSTALLER_SIGNING_IDENTITY}" output/default-browser.pkg

# notarize the package
echo "Notarizing the package"
$XCODE_NOTARY_PATH store-credentials --apple-id "opensource@macadmins.io" --team-id "T4SK8ZXCXG" --password "$2" defaultbrowser

# Notarize default-browser package
$XCODE_NOTARY_PATH submit "output/default-browser.pkg" --keychain-profile "defaultbrowser" --wait
$XCODE_STAPLER_PATH staple "output/default-browser.pkg"