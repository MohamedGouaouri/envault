package cmd

import (
	"os"
	"strings"

	"github.com/MohamedGouaouri/envault/config"
	"github.com/MohamedGouaouri/envault/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "envault",
	Short:             "envault CLI is used to inject environment variables into any process from HashiCorp Vault",
	Long:              `envault is a simple, end-to-end encrypted service that enables teams to sync and manage their environment variables across their development life cycle.`,
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	Version:           util.CLI_VERSION,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "log level (trace, debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().StringVar(&config.VAULT_ADDR, "domain", util.VAULT_ADDR_DEFAULT, "Point the CLI to your own backend [can also set via environment variable name: VAULT_ADDR]")

	// if config.VAULT_ADDR is set to the default value, check if VAULT_ADDR is set in the environment
	// this is used to allow overrides of the default value
	if !rootCmd.Flag("domain").Changed {
		if envVaultAddr, ok := os.LookupEnv("VAULT_ADDR"); ok {
			config.VAULT_ADDR = envVaultAddr
		}
	}

}

func initLog() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	ll, err := rootCmd.Flags().GetString("log-level")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	switch strings.ToLower(ll) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "err", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
