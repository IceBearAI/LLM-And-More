package service

import (
	"github.com/spf13/cobra"
)

var (
	jobCmd = &cobra.Command{
		Use:               "job command <args> [flags]",
		Short:             "任务命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
可用的配置类型：
[finetuning]

aigc-server job -h
`,
	}
)
