package wyago

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "wyago"}
	rootCmd.AddCommand(Commands...)
	rootCmd.Execute()
}
