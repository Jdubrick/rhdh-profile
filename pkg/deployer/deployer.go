package deployer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Jdubrick/rhdh-profile/pkg/config"
	"github.com/Jdubrick/rhdh-profile/pkg/kustomize"
)

type Deployer struct {
	config     *config.Config
	kubeClient kubernetes.Interface
	kustomizer *kustomize.Kustomizer
}

func New(cfg *config.Config) (*Deployer, error) {

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	kustomizer, err := kustomize.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kustomizer: %w", err)
	}

	return &Deployer{
		config:     cfg,
		kubeClient: kubeClient,
		kustomizer: kustomizer,
	}, nil
}

func (d *Deployer) DeployOperator() error {
	if d.config.Verbose {
		fmt.Println("Deploying RHDH operator...")
	}

	profilePath := d.config.ProfilePath
	if !filepath.IsAbs(profilePath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		profilePath = filepath.Join(wd, profilePath)
	}

	if d.config.DryRun {
		fmt.Printf("Dry run - would run: kubectl apply -k %s\n", profilePath)
		return nil
	}

	if err := d.applyKustomize(profilePath); err != nil {
		return fmt.Errorf("failed to apply operator manifests: %w", err)
	}

	if d.config.Verbose {
		fmt.Println("RHDH operator deployed successfully!")
	}

	return nil
}

func (d *Deployer) DeployPresets() error {
	if d.config.Verbose {
		fmt.Println("Deploying RHDH presets...")
	}

	if err := d.updateRHDHSecrets(); err != nil {
		return fmt.Errorf("failed to update RHDH secrets: %w", err)
	}

	presetsPath := filepath.Join(d.config.PresetsPath, "rhdh-complete")
	if !filepath.IsAbs(presetsPath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		presetsPath = filepath.Join(wd, presetsPath)
	}

	if d.config.DryRun {
		fmt.Printf("Dry run - would run: kubectl apply -k %s\n", presetsPath)
		return nil
	}

	if err := d.applyKustomize(presetsPath); err != nil {
		return fmt.Errorf("failed to apply presets manifests: %w", err)
	}

	if d.config.Verbose {
		fmt.Println("RHDH presets deployed successfully!")
	}

	return nil
}

func (d *Deployer) updateRHDHSecrets() error {
	if len(d.config.EnvVars) == 0 {
		if d.config.Verbose {
			fmt.Println("No environment variables found, skipping secrets update")
		}
		return nil
	}

	secretsPath := filepath.Join("presets", "rhdh-complete", "rhdh", "rhdh-secrets.yaml")

	if d.config.Verbose {
		fmt.Printf("Updating RHDH secrets file: %s\n", secretsPath)
	}

	return d.kustomizer.UpdateSecretsFile(secretsPath, d.config.EnvVars)
}

func (d *Deployer) applyKustomize(kustomizePath string) error {
	args := []string{"apply", "-k", kustomizePath}

	cmd := exec.Command("kubectl", args...)
	if d.config.Verbose {
		fmt.Printf("Running: kubectl %s\n", strings.Join(args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kubectl apply -k failed: %w", err)
	}

	return nil
}
