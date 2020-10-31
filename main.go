package main

import (
	"wyago/cmd"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "wyago"}
	rootCmd.AddCommand(cmd.Commands...)
	rootCmd.Execute()
}
