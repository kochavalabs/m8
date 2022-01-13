package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kochavalabs/mazzaroth-cli/internal/cfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func Execute() error {
	// root command entry to application
	rootCmd := &cobra.Command{
		Use:   "m8",
		Short: "mazzaroth command line interface",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			dir, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			if _, err := os.Stat(dir + cfgDir + cfgName); errors.Is(err, os.ErrNotExist) {
				return errors.New(err.Error() + ", please run m8 cfg init to create the cfg")
			}

			cfg, err := cfg.FromFile(dir + cfgDir + cfgName)
			if err != nil {
				return err
			}

			fmt.Println(cfg)

			viper.Set("cfg", cfg)
			// Bind Cobra flags with viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Environment variables are expected to be ALL CAPS
			viper.AutomaticEnv()
			viper.SetEnvPrefix("m8")
			return nil
		},
	}

	rootCmd.AddCommand(blockCmdChain())
	rootCmd.AddCommand(channelCmdChain())
	rootCmd.AddCommand(configurationCmdChain())
	rootCmd.AddCommand(receiptCmdChain())
	rootCmd.AddCommand(transactionCmdChain())
	rootCmd.PersistentFlags().String(cfgPath, "$HOME/.m8/cfg.yaml", "location of the mazzaroth config file")
	rootCmd.PersistentFlags().String(channelId, "", "defaults to the active channel id in the cfg")
	rootCmd.PersistentFlags().String(address, "", "defaults to active channel address in the cfg")

	ctx, cancel := context.WithCancel(context.Background())
	errGrp, errctx := errgroup.WithContext(ctx)
	errGrp.Go(func() error {
		defer cancel()
		if err := rootCmd.ExecuteContext(errctx); err != nil {
			return err
		}
		return nil
	})
	return errGrp.Wait()
}
