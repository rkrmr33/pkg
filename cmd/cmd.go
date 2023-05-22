package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// VersionInfo contains version information
type VersionInfo struct {
	// Bin binary name
	Bin string
	// Version version name (e.g. 1.99.101)
	Version string
	// BuildDate build date
	BuildDate string
	// GitCommit git commit SHA
	GitCommit string
}

// NewVersionCmd returns a command that prints the version information
func NewVersionCmd(versionInfo VersionInfo) *cobra.Command {
	var (
		short bool
	)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints binary version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printVersionInfo(short, versionInfo)
		},
	}

	cmd.Flags().BoolVar(&short, "short", false, "show only binary version")

	return cmd
}

func printVersionInfo(short bool, versionInfo VersionInfo) error {
	fmt.Printf("%s: %s\n", versionInfo.Bin, versionInfo.Version)

	if short {
		return nil
	}

	fmt.Printf("  BuildDate: %s\n", versionInfo.BuildDate)
	fmt.Printf("  GitCommit: %s\n", versionInfo.GitCommit)
	fmt.Printf("  GoVersion: %s\n", runtime.Version())
	fmt.Printf("  CompilerVersion: %s\n", runtime.Compiler)
	fmt.Printf("  Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	return nil
}