package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version number of helmtrans",
    Long:  `All software has versions. This is helmtrans's`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("helmtrans version is v0.1.0")
    },
}