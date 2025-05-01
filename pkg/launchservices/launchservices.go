package launchservices

import (
	"fmt"
	"os"

	"github.com/macadmins/default-browser/pkg/client"
)

const lsregister = "/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister"

func ModifyLS(c client.Client, identifier string, noRebuildLaunchServices bool) error {
	plist, err := GetPlist(c.PlistLocation)
	if err != nil {
		return err
	}

	plist.CleanHandlers()

	plist.AddLSHandlers([]LSHandler{
		{LSHandlerContentType: "public.url", LSHandlerPreferredVersions: map[string]string{"LSHandlerRoleViewer": "-"}, LSHandlerRoleViewer: identifier},
		{LSHandlerContentType: "public.xhtml", LSHandlerPreferredVersions: map[string]string{"LSHandlerRoleAll": "-"}, LSHandlerRoleAll: identifier},
		{LSHandlerContentType: "public.html", LSHandlerPreferredVersions: map[string]string{"LSHandlerRoleAll": "-"}, LSHandlerRoleAll: identifier},
		{LSHandlerPreferredVersions: map[string]string{"LSHandlerRoleAll": "-"}, LSHandlerRoleAll: identifier, LSHandlerURLScheme: "https"},
		{LSHandlerPreferredVersions: map[string]string{"LSHandlerRoleAll": "-"}, LSHandlerRoleAll: identifier, LSHandlerURLScheme: "http"},
	})

	err = WritePlist(c.PlistLocation, plist)
	if err != nil {
		return err
	}

	lsRegisterExists, err := lsRegisterExists(lsregister)
	if err != nil {
		return err
	}

	if !lsRegisterExists {
		fmt.Printf("lsregister does not exist at %s. You should restart the device to rebuild launchservices", lsregister)
		return nil
	}

	err = rebuildLaunchServices(c, noRebuildLaunchServices)
	if err != nil {
		return err
	}

	err = killlsd(c)
	if err != nil {
		return err
	}

	return nil
}

func rebuildLaunchServices(c client.Client, noRebuildLaunchServices bool) error {
	if !noRebuildLaunchServices {
		_, err := c.Runner.RunCmd(lsregister, "-gc", "-R", "-all", "user,system,local,network")
		if err != nil {
			return err
		}
	}
	return nil
}

func killlsd(c client.Client) error {
	_, err := c.Runner.RunCmd("/usr/bin/killall", "lsd")
	if err != nil {
		return err
	}
	return nil
}

func lsRegisterExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
