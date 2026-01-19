package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 这些变量将在编译时由 GoReleaser 通过 -ldflags 注入
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("aicommits %s (commit: %s, built at: %s)\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
