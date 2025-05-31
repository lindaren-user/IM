package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "启动项目",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, IM!")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
