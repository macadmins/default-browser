package launchservices

import (
	"os"
	"path/filepath"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed binary_data.plist
var binaryplist []byte

func TestReadPlist(t *testing.T) {
	// Create a temporary file with fake plist data
	tmpfile, err := os.CreateTemp("", "test.plist")
	assert.NoError(t, err, "TempFile should not return an error")
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(binaryplist))
	assert.NoError(t, err, "Write should not return an error")
	err = tmpfile.Close()
	assert.NoError(t, err, "Close should not return an error")

	// Read the plist file
	p, err := ReadPlist(tmpfile.Name())
	assert.NoError(t, err, "ReadPlist should not return an error")

	// Verify the contents of the plist
	assert.Len(t, p.LSHandlers, 1, "LSHandlers should have one entry")
	assert.Equal(t, "public.html", p.LSHandlers[0].LSHandlerContentType, "LSHandlerContentType should match")
	assert.Equal(t, "com.apple.safari", p.LSHandlers[0].LSHandlerRoleAll, "LSHandlerRoleAll should match")
}

func TestWritePlist(t *testing.T) {
	// Create a temporary file to write the plist data
	tmpfile, err := os.CreateTemp("", "test.plist")
	assert.NoError(t, err, "TempFile should not return an error")
	defer os.Remove(tmpfile.Name())

	// Create a Plist object with fake data
	p := Plist{
		LSHandlers: []LSHandler{
			{
				LSHandlerContentType: "public.html",
				LSHandlerRoleAll:     "com.apple.safari",
			},
		},
	}

	// Write the plist data to the temporary file
	err = WritePlist(tmpfile.Name(), p)
	assert.NoError(t, err, "WritePlist should not return an error")

	// Read the plist file to verify its contents
	readPlist, err := ReadPlist(tmpfile.Name())
	assert.NoError(t, err, "ReadPlist should not return an error")

	// Verify the contents of the plist
	assert.Len(t, readPlist.LSHandlers, 1, "LSHandlers should have one entry")
	assert.Equal(t, "public.html", readPlist.LSHandlers[0].LSHandlerContentType, "LSHandlerContentType should match")
	assert.Equal(t, "com.apple.safari", readPlist.LSHandlers[0].LSHandlerRoleAll, "LSHandlerRoleAll should match")
}

func TestPlistExists(t *testing.T) {
	// Create a temporary file to test existence
	tmpfile, err := os.CreateTemp("", "test.plist")
	assert.NoError(t, err, "TempFile should not return an error")
	defer os.Remove(tmpfile.Name())

	// Test that the file exists
	exists, err := plistExists(tmpfile.Name())
	assert.NoError(t, err, "PlistExists should not return an error")
	assert.True(t, exists, "PlistExists should return true for an existing file")

	// Remove the temporary file
	err = os.Remove(tmpfile.Name())
	assert.NoError(t, err, "Remove should not return an error")

	// Test that the file does not exist
	exists, err = plistExists(tmpfile.Name())
	assert.NoError(t, err, "PlistExists should not return an error for a non-existing file")
	assert.False(t, exists, "PlistExists should return false for a non-existing file")
}

func TestGetPlist(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err, "TempDir should not return an error")
	defer os.RemoveAll(tmpDir)

	// Define the path for the temporary plist file
	tmpFilePath := filepath.Join(tmpDir, "test.plist")

	// Create a Plist object with fake data
	expectedPlist := Plist{
		LSHandlers: []LSHandler{
			{
				LSHandlerContentType: "public.html",
				LSHandlerRoleAll:     "com.apple.safari",
			},
		},
	}

	// Write the fake plist data to the temporary file
	err = WritePlist(tmpFilePath, expectedPlist)
	assert.NoError(t, err, "WritePlist should not return an error")

	// Test GetPlist function
	p, err := GetPlist(tmpFilePath)
	assert.NoError(t, err, "GetPlist should not return an error")
	assert.Equal(t, expectedPlist, p, "GetPlist should return the correct Plist object")

	// Test GetPlist function with a non-existing file
	nonExistentFilePath := filepath.Join(tmpDir, "nonexistent.plist")
	p, err = GetPlist(nonExistentFilePath)
	assert.NoError(t, err, "GetPlist should not return an error for a non-existing file")
	assert.Empty(t, p.LSHandlers, "GetPlist should return an empty Plist object for a non-existing file")
}
