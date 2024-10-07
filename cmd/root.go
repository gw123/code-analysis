package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd 定义了主命令
var rootCmd = &cobra.Command{
	Use:   "code-analyzer",
	Short: "Analyze and summarize code using AI",
}

// Execute 启动命令行工具
func Execute() error {
	return rootCmd.Execute()
}
