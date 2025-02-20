package launchservices

import (
	"os"
	"path/filepath"

	"github.com/micromdm/plist"
)

type LSHandler struct {
	LSHandlerContentType       string            `plist:"LSHandlerContentType,omitempty"`
	LSHandlerPreferredVersions map[string]string `plist:"LSHandlerPreferredVersions,omitempty"`
	LSHandlerRoleViewer        string            `plist:"LSHandlerRoleViewer,omitempty"`
	LSHandlerRoleAll           string            `plist:"LSHandlerRoleAll,omitempty"`
	LSHandlerURLScheme         string            `plist:"LSHandlerURLScheme,omitempty"`
}

type Plist struct {
	LSHandlers []LSHandler `plist:"LSHandlers"`
}

func GetPlist(path string) (Plist, error) {
	var p Plist
	exists, err := plistExists(path)
	if err != nil {
		return p, err
	}

	if exists {
		p, err = ReadPlist(path)
		if err != nil {
			return p, err
		}
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return p, err
	}
	return p, nil
}

func plistExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func ReadPlist(path string) (Plist, error) {
	var p Plist
	_, err := os.Stat(path)
	if err != nil {
		return p, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}

	err = plist.Unmarshal(data, &p)
	if err != nil {
		return p, err
	}

	return p, nil
}

func WritePlist(path string, p Plist) error {
	plistFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer plistFile.Close()

	err = plist.NewEncoder(plistFile).Encode(p)
	if err != nil {
		return err
	}

	return nil
}
