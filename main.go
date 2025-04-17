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
	var noRebuildLaunchServices bool

	var rootCmd = &cobra.Command{
		Use:   "default-browser",
		Short: "A cli tool to set the default browser on macOS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setDefault(identifier, noRebuildLaunchServices)
		},
	}

	rootCmd.Flags().StringVar(&identifier, "identifier", "com.google.chrome", "An identifier for the application")
	rootCmd.Flags().BoolVar(&noRebuildLaunchServices, "no-rebuild-launchservices", false, "Do not rebuild launch services. Only use if you are experiencing issues with System Settings not displaying correctly after a reboot.")

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("default-browser version {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setDefault(identifier string, noRebuildLaunchServices bool) error {
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

	err = launchservices.ModifyLS(c, identifier, noRebuildLaunchServices)
	if err != nil {
		return err
	}
	return nil
}
