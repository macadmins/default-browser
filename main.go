package main

import (
	"fmt"
	"os"

	"github.com/macadmins/default-browser/pkg/client"
	"github.com/macadmins/default-browser/pkg/launchservices"
	"github.com/spf13/cobra"
)

var version string

func main() {
	var identifier string
	var noRescanLaunchServices bool

	var rootCmd = &cobra.Command{
		Use:   "default-browser",
		Short: "A cli tool to set the default browser on macOS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setDefault(identifier, noRescanLaunchServices)
		},
	}

	rootCmd.Flags().StringVar(&identifier, "identifier", "com.google.chrome", "An identifier for the application")
	rootCmd.Flags().BoolVar(&noRescanLaunchServices, "no-rescan-launchservices", false, "Do not rescan launch services. Only use if you are experiencing issues with System Settings not displaying correctly after a reboot.")
	rootCmd.Flags().BoolVar(&noRescanLaunchServices, "no-rebuild-launchservices", false, "Legacy: same as --no-rescan-launchservices")

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("default-browser version {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setDefault(identifier string, noRescanLaunchServices bool) error {
	if identifier == "" {
		return fmt.Errorf("identifier cannot be empty")
	}

	// Todo: actually run as the logged in user if run as root. For now just bail
	if os.Geteuid() == 0 {
		return fmt.Errorf("this tool must be run as the logged in user")
	}

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	err = launchservices.ModifyLS(c, identifier, noRescanLaunchServices)
	if err != nil {
		return err
	}
	return nil
}
