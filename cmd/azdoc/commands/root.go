package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version   string
	gitCommit string
	buildDate string
	cfgFile   string
	logLevel  string
	noAnsi    bool
	quiet     bool
)

var rootCmd = &cobra.Command{
	Use:   "azdoc",
	Short: "Azure documentation generator",
	Long: `azdoc - Generate comprehensive documentation for Azure subscriptions.

Creates inventory tables, network topology diagrams, routing views, and more.
Supports optional LLM-powered explanations with deterministic outputs.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./azdoc.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug|info|warn|error)")
	rootCmd.PersistentFlags().BoolVar(&noAnsi, "no-ansi", false, "disable ANSI colors in output")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress all non-error output")

	// Bind flags to viper
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("no-ansi", rootCmd.PersistentFlags().Lookup("no-ansi"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("azdoc")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.azdoc")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("AZDOC")

	// Ignore error if config file not found
	if err := viper.ReadInConfig(); err == nil {
		// Config file found and loaded
	}
}

func SetVersion(v, commit, date string) {
	version = v
	gitCommit = commit
	buildDate = date
}

func Execute() error {
	return rootCmd.Execute()
}

func GetVersion() string {
	return fmt.Sprintf("azdoc version %s (commit: %s, built: %s)", version, gitCommit, buildDate)
}
