package cmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/m8/internal/manifest"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func channelCmdChain() *cobra.Command {

	channelRootCmd := &cobra.Command{
		Use:   "channel",
		Short: "interact with channel endpoints on a mazzaroth gateway node",
	}
	channelRootCmd.PersistentFlags().String(channelId, "", "defaults to the active channel id in the cfg")
	channelRootCmd.PersistentFlags().String(channelAddress, "", "defaults to active channel address in the cfg")

	channelAbiCmd := &cobra.Command{
		Use:   "abi",
		Short: "return the application binary interface for a channel",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			abi, err := client.ChannelAbi(cmd.Context(), viper.GetString(channelId))
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(abi, "", "\t")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}

	channelGenCmd := &cobra.Command{
		Use:   "gen",
		Short: "generate a mazzaroth channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate Key
			pub, priv, err := crypto.GenerateEd25519KeyPair()
			if err != nil {
				return err
			}
			fmt.Println("channel address:", crypto.ToHex(pub))
			fmt.Println("channel private key:", crypto.ToHex(priv))
			fmt.Println("channel Id:", crypto.ToHex(pub))
			// TODO
			// self signed cert for channel
			// mazzaroth.io cert generate for channel
			return nil
		},
	}

	channelDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a channel contract",
		RunE: func(cmd *cobra.Command, args []string) error {

			sender, err := xdr.IDFromHexString(viper.GetString(publicKey))
			if err != nil {
				return err
			}

			channelId, err := xdr.IDFromHexString(viper.GetString(channelId))
			if err != nil {
				return err
			}

			pk, err := crypto.FromHex(viper.GetString(privateKey))
			if err != nil {
				return err
			}

			tx, err := mazzaroth.Transaction(sender, channelId).Contract(mazzaroth.GenerateNonce(), defaultBlockExpirationNumber).Delete().Sign(pk)
			if err != nil {
				return err
			}

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			id, _, err := client.TransactionSubmit(cmd.Context(), tx)
			if err != nil {
				return err
			}

			fmt.Println("transaction id:", hex.EncodeToString(id[:]))
			return nil
		},
	}

	channelDeployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy a channel contract to mazzaroth nodes from a given manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := viper.GetString(deploymentManifest)
			// check if file is in default path
			if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
				return errors.New("unable to locate deployment manifest")
			}

			pk, err := crypto.FromHex(viper.GetString(privateKey))
			if err != nil {
				return err
			}

			manifests, err := manifest.FromFile(manifestPath, "deployment")
			if err != nil {
				return err
			}

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			if err := manifest.ExecuteDeployments(cmd.Context(), manifests, client, viper.GetString(publicKey), pk); err != nil {
				return err
			}
			return nil
		},
	}
	channelDeployCmd.Flags().String(deploymentManifest, defaultDeploymentManifestPath, "location of mazzaroth channel deployment manifest")

	channelPauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "pause/unpause a channel contract",
		RunE: func(cmd *cobra.Command, args []string) error {

			pk, err := crypto.FromHex(viper.GetString(privateKey))
			if err != nil {
				return err
			}

			sender, err := xdr.IDFromHexString(viper.GetString(publicKey))
			if err != nil {
				return err
			}

			channelId, err := xdr.IDFromHexString(viper.GetString(channelId))
			if err != nil {
				return err
			}

			tx, err := mazzaroth.Transaction(sender, channelId).Contract(mazzaroth.GenerateNonce(), defaultBlockExpirationNumber).
				Pause(viper.GetBool(pause)).
				Sign(pk)
			if err != nil {
				return err
			}

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			id, _, err := client.TransactionSubmit(cmd.Context(), tx)
			if err != nil {
				return err
			}

			fmt.Println("transaction id:", id)
			return nil
		},
	}
	channelPauseCmd.Flags().Bool(pause, false, "pause transactions from being sent")

	channelTestCmd := &cobra.Command{
		Use:   "test",
		Short: "test channel contracts on mazzaroth nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := viper.GetString(testManifest)
			// check if file is in default path
			if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
				return errors.New("unable to locate test manifest")
			}

			manifests, err := manifest.FromFile(manifestPath, "test")
			if err != nil {
				return err
			}

			pk, err := crypto.FromHex(viper.GetString(privateKey))
			if err != nil {
				return err
			}

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			if err := manifest.ExecuteTests(cmd.Context(), manifests, client, viper.GetString(publicKey), pk); err != nil {
				return err
			}

			return nil
		},
	}
	channelTestCmd.Flags().String(testManifest, defaultTestManifestPath, "location of mazzaroth channel test manifest")

	channelRootCmd.AddCommand(
		channelAbiCmd,
		channelGenCmd,
		channelDeleteCmd,
		channelDeployCmd,
		channelPauseCmd,
		channelTestCmd,
		blockCmdChain(),
		transactionCmdChain(),
		receiptCmdChain(),
	)
	return channelRootCmd
}
