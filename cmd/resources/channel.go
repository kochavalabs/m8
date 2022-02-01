package resources

import (
	"errors"
	"fmt"
	"os"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/m8/cmd/verbs"
	"github.com/kochavalabs/m8/internal/manifest"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ChannelCmdChain() *cobra.Command {

	channelRootCmd := &cobra.Command{
		Use:   "channel",
		Short: "interact with channel endpoints on a mazzaroth gateway node",
	}
	channelRootCmd.PersistentFlags().String(channelId, "", "defaults to the active channel id in the cfg")
	channelRootCmd.PersistentFlags().String(channelAddress, "", "defaults to active channel address in the cfg")

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
		channelGenCmd,
		channelDeployCmd,
		channelTestCmd,
		verbs.Lookup("channel"),
		verbs.Exec(),
		verbs.Pause())
	return channelRootCmd
}
