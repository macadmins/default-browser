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

	var rootCmd = &cobra.Command{
		Use:   "default-browser",
		Short: "A cli tool to set the default browser on macOS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setDefault(identifier)
		},
	}

	rootCmd.Flags().StringVar(&identifier, "identifier", "com.google.chrome", "An identifier for the application")

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("default-browser version {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setDefault(identifier string) error {
	if identifier == "" {
		return fmt.Errorf("identifier cannot be empty")
	}

	c, err := client.NewClient()
	if err != nil {
		return err
	}

	err = launchservices.ModifyLS(c, identifier)
	if err != nil {
		return err
	}
	return nil
}
