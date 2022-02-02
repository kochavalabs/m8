package channel

import (
	"errors"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/m8/internal/manifest"
	"github.com/kochavalabs/m8/internal/tui"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	function   = `fn`
	arguments  = `args`
	privateKey = `private-key`
	publicKey  = `public-key`

	maxBlockExpirationRange       = 10
	deploymentManifest            = `deployment-manifest`
	testManifest                  = `test-manifest`
	defaultDeploymentManifestPath = `./m8/deployment.yaml`
	defaultTestManifestPath       = `./m8/test.yaml`
)

func exec() *cobra.Command {
	exec := &cobra.Command{
		Use:   "exec",
		Short: "preform executions against a mazzaroth channel",
	}
	exec.AddCommand(execTx(), execDeployment(), execTest())
	return exec
}

func execTx() *cobra.Command {
	execTx := &cobra.Command{
		Use:   "tx",
		Short: "execute functions on a mazzaroth channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}
			pk, err := crypto.FromHex(viper.GetString(privateKey))
			if err != nil {
				return err
			}
			sender, err := xdr.IDFromHexString(viper.GetString(publicKey))
			if err != nil {
				return err
			}

			cId, err := xdr.IDFromHexString(viper.GetString(channelId))
			if err != nil {
				return err
			}

			xdrArgs := make([]xdr.Argument, 0, 0)
			values := viper.GetStringSlice(arguments)
			for _, a := range values {
				xdrArgs = append(xdrArgs, xdr.Argument(a))
			}

			blockHeight, err := client.BlockHeight(cmd.Context(), viper.GetString(channelId))
			if err != nil {
				return err
			}

			tx, err := mazzaroth.Transaction(sender, cId).
				Call(mazzaroth.GenerateNonce(), blockHeight.Height+maxBlockExpirationRange).
				Function(viper.GetString(function)).
				Arguments(xdrArgs...).
				Sign(pk)
			if err != nil {
				return err
			}

			txCmd := tui.TxCall(cmd.Context(), client, tx)
			txModel := tui.NewTxModel(txCmd)

			if err := tea.NewProgram(txModel).Start(); err != nil {
				return err
			}
			return nil
		},
	}
	execTx.Flags().String(function, "", "the function to be called")
	execTx.MarkFlagRequired(function)

	execTx.Flags().StringArray(arguments, []string{""}, "the args to pass within the function")
	return execTx
}

func execDeployment() *cobra.Command {
	execDeployment := &cobra.Command{
		Use:   "deployment",
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
	execDeployment.Flags().String(deploymentManifest, defaultDeploymentManifestPath, "location of mazzaroth channel deployment manifest")
	return execDeployment
}

func execTest() *cobra.Command {
	execTest := &cobra.Command{
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
	execTest.Flags().String(testManifest, defaultTestManifestPath, "location of mazzaroth channel test manifest")
	return execTest
}
