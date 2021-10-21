package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// Environment variables are expected to be ALL CAPS
	// ENV CONST VARS
	// log level debug|info|warn|error|dpanic|panic|fatal
	logLevel = "log-level"
	cfgPath  = "cfg-path"
)

//rootCmd - is the main command and runs a series of "setup" functions that are passed to child cmds
var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "Mazzaroth Cli",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Bind Cobra flags with viper
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		// Environment variables are expected to be ALL CAPS
		viper.AutomaticEnv()
		viper.SetEnvPrefix("mazzaroth_cli")
		// create zap logger
		zlog, err := initLogger(viper.GetString(logLevel))
		if err != nil {
			return err
		}
		// Set the logger in viper registery
		viper.Set("logger", zlog)
		return nil
	},
}

// initLogger - initalizes a zap production logger
func initLogger(loglevel string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	switch strings.ToLower(loglevel) {
	case "debug":
		config.Level.SetLevel(zap.DebugLevel)
	case "warn":
		config.Level.SetLevel(zap.WarnLevel)
	case "error":
		config.Level.SetLevel(zap.ErrorLevel)
	case "dpanic":
		config.Level.SetLevel(zap.DPanicLevel)
	case "panic":
		config.Level.SetLevel(zap.PanicLevel)
	case "fatal":
		config.Level.SetLevel(zap.FatalLevel)
	default:
		// Default info level
		config.Level.SetLevel(zap.InfoLevel)
	}
	// Create the zap logger from config
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func Execute() error {
	// cmd chain
	// root
	//// setup
	//// channel
	////// list
	////// deploy
	////// connect
	////// contract
	////// abi
	////// functions
	//// transaction
	////// lookup
	////// list
	////// call
	////// update
	////// receipt
	//////// lookup
	//// block
	////// lookup
	////// list

	rootCmd.AddCommand(configure())
	rootCmd.PersistentFlags().String(logLevel, "info", "Log level :: debug|info|warn|error|dpanic|panic|fatal")
	rootCmd.PersistentFlags().String(cfgPath, "$HOME/.mazzaroth-cli", "Location of the mazzaroth cli file")
	return rootCmd.Execute()
}
