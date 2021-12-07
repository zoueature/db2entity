/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "db2entity",
	Short: "同步数据库表到go实体",
	Long:  `同步数据库表到go实体`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := initDB(
			cmd.Flag("host").Value.String(),
			cmd.Flag("port").Value.String(),
			cmd.Flag("username").Value.String(),
			cmd.Flag("password").Value.String(),
		)
		if err != nil {
			panic(err.Error())
		}
		err = synTable(cmd.Flag("database").Value.String(), cmd.Flag("destination").Value.String(), cmd.Flag("package").Value.String(), cmd.Flag("prefix").Value.String())
		if err != nil {
			println(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.db2struct.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("help", "", false, "usage help")
	rootCmd.Flags().StringP("host", "h", "127.0.0.1", "数据库host")
	rootCmd.Flags().StringP("port", "P", "3306", "数据库端口号")
	rootCmd.Flags().StringP("password", "p", "", "数据库密码")
	rootCmd.Flags().StringP("username", "u", "root", "数据库用户名")
	rootCmd.Flags().StringP("database", "d", "", "数据库名")
	rootCmd.Flags().StringP("destination", "", ".", "生成是实体目录")
	rootCmd.Flags().StringP("package", "", "entity", "包名")
	rootCmd.Flags().StringP("prefix", "", "", "需要去除的表前缀")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".app" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".db2struct")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
