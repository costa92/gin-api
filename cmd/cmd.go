package cmd

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "api is web",
	Long:  "A Fast and at http://test",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Info().Msgf("homedir dir %s:", err)
			os.Exit(0)
		}
		viper.AddConfigPath(home)
		viper.SetConfigFile(".cobra")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Info().Msgf("Config read config:", err)
		os.Exit(0)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Msgf("exec fail:", err)
		os.Exit(0)
	}
}
