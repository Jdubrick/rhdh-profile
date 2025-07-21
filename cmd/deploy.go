package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Jdubrick/rhdh-profile/pkg/config"
	"github.com/Jdubrick/rhdh-profile/pkg/deployer"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy RHDH components",
	Long: `Deploy RHDH components: either the operator or the presets.
The deploy command supports deploying individual components.`,
}

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Deploy the RHDH operator",
	Long: `Deploy the RHDH operator using the Kustomize manifests in the profile directory.
This will install the Red Hat Developer Hub operator in the specified namespace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig(envFile)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		d, err := deployer.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to create deployer: %w", err)
		}

		if viper.GetBool("verbose") {
			fmt.Println("Deploying RHDH operator...")
		}

		return d.DeployOperator()
	},
}

var presetsCmd = &cobra.Command{
	Use:   "presets",
	Short: "Deploy the RHDH presets",
	Long: `Deploy the RHDH presets using the Kustomize manifests in the presets directory.
If an .env file is provided, it will update the rhdh-secrets.yaml file with values from the .env file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig(envFile)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		d, err := deployer.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to create deployer: %w", err)
		}

		if viper.GetBool("verbose") {
			fmt.Println("Deploying RHDH preset...")
		}

		return d.DeployPresets()
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.AddCommand(operatorCmd)
	deployCmd.AddCommand(presetsCmd)

	deployCmd.PersistentFlags().Int("timeout", 600, "timeout in seconds for deployment operations")
	viper.BindPFlag("timeout", deployCmd.PersistentFlags().Lookup("timeout"))
}
