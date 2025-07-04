package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "im",
	Short: "即时聊天系统",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("出错啦：", err)
		os.Exit(1)
	}
}
