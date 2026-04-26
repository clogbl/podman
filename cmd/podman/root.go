package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// rootCmd is the base command for the podman CLI.
	rootCmd = &cobra.Command{
		Use:   "podman",
		Short: "Manage pods, containers and images",
		Long: `Podman is a tool for managing pods, containers, and container images.

It provides a Docker-compatible command line interface for managing containers
without requiring a daemon. Podman can run containers as root or in rootless mode.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Global flags
	logLevel   string
	connection string
	identity   string
	noOut      bool
)

func init() {
	// Persistent flags are available to all subcommands.
	rootCmd.PersistentFlags().StringVar(
		&logLevel,
		"log-level",
		"warn",
		`Log messages above specified level (trace, debug, info, warn, error, fatal, panic)`,
	)
	rootCmd.PersistentFlags().StringVarP(
		&connection,
		"connection", "c",
		"",
		"Connection to use for remote Podman service",
	)
	rootCmd.PersistentFlags().StringVar(
		&identity,
		"identity",
		"",
		"Path to SSH identity file, (CONTAINER_SSHKEY)",
	)
	rootCmd.PersistentFlags().BoolVar(
		&noOut,
		"noout",
		false,
		"Do not output to stdout",
	)
}

// Execute runs the root command and handles top-level errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
