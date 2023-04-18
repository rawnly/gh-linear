package main

import (
	"context"

	cmd "github.com/Rawnly/gh-linear/cmd"
	cfg "github.com/Rawnly/gh-linear/config"
	"github.com/Rawnly/gh-linear/linear"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// setup viper
	viper.SetConfigType("json")
	viper.SetConfigName("gh-linear")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("$HOME/.gh-linear")
	viper.SetEnvPrefix("gh_linear")

	// read .env or from environment
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Debug("No settings file found")

			// Create file if not exists
			if err := viper.SafeWriteConfig(); err != nil {
				cobra.CheckErr(err)
			}

			if _, err := cfg.LoadDefaults(); err != nil {
				cobra.CheckErr(err)
			}
		} else {
			logrus.Fatal("Error reading settings file:", err)
		}
	}

	logrus.Debug("Config loaded")
	logrus.Debugf("Using %s", viper.ConfigFileUsed())
}

func main() {
	linearClient := linear.NewClient()
	logrus.SetLevel(logrus.DebugLevel)

	config, err := cfg.Load()
	if err != nil {
		logrus.Fatal(err)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "linear", linearClient)
	ctx = context.WithValue(ctx, "config", &config)

	cmd.Execute(ctx)
}
