package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra" // 开源很好用的第三方cli组件
)

var (
	vers bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	// Use为命令名名称
	Use: "restful-api",
	// Short和Long为命令说明
	Short: "restful-api CRUD Demo",
	Long:  `restful-api ...`,
	// RunE为输入Use命令的处理函数
	RunE: func(cmd *cobra.Command, args []string) error {
		// 打印版本信息
		if vers {
			fmt.Println("0.0.1")
			return nil
		}
		// 执行根命令未加参数时的报错信息，会将所有根选项打印出来
		return errors.New("no flags find")
	},
}

func Execute() {
	// 调用RootCmd命令
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 持久化参数，设置变量命令行参数，格式为：变量、长参数、短参数、默认值、命令说明
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the restful-api version")
}
