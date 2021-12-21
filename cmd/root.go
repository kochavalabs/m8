package cmd

import (
	"context"

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
	rootCmd.AddCommand(blockCmdChain())
	rootCmd.AddCommand(channelCmdChain())
	rootCmd.AddCommand(configure())
	rootCmd.AddCommand(receiptCmdChain())
	rootCmd.PersistentFlags().String(cfgPath, "$HOME/.m8", "location of the mazzaroth config file")
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
