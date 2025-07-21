package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var envFile string

var rootCmd = &cobra.Command{
	Use:   "rhdh-cli",
	Short: "A CLI tool for deploying RHDH (Red Hat Developer Hub) components",
	Long: `rhdh-cli is a command line tool that wraps Kustomize-based Kubernetes deployments
for Red Hat Developer Hub. It provides commands to deploy the RHDH operator and
complete RHDH installations with automatic resource monitoring and configuration.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rhdh-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&envFile, "env-file", ".env", "environment file to load variables from")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "perform a dry run without actually deploying")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rhdh-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
