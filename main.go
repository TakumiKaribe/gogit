package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	print := &cobra.Command{
		Use:   "print use",
		Short: "print short",
		Long:  "print long",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: ", strings.Join(args, " "))
		},
	}
	rootCmd := &cobra.Command{Use: "wyago"}
	rootCmd.AddCommand(print)
	rootCmd.Execute()
}
