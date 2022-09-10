package cmd

import (
	"github.com/costa92/go-web/internal/logger"
	"os"

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
		viper.AddConfigPath(".")
		viper.SetConfigFile("config.yaml")
	}
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorw("viper ReadInConfig", "err", err)
		os.Exit(0)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(0)
	}
}
