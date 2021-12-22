package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/kochavalabs/mazzaroth-cli/internal/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func deployCmdChain() *cobra.Command {

	deployRootCmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy channel contracts to mazzaroth nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := viper.GetString(deploymentManifestPath)
			// check if file is in default path
			if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
				return errors.New("unable to locate deployment manifest")
			}
			manifestFile, err := ioutil.ReadFile(manifestPath)
			if err != nil {
				return err
			}

			manifests := make([]*manifest.Manifest, 0, 0)
			r := bytes.NewReader(manifestFile)

			decoder := yaml.NewDecoder(r)
			for {
				manifest := &manifest.Manifest{}
				if err := decoder.Decode(manifest); err != nil {
					if err != io.EOF {
						return os.ErrClosed
					}
					break
				}
				if manifest.Type == "deployment" {
					manifests = append(manifests, manifest)
				}
			}
			fmt.Println(manifests[0].Channel.AbiFile)
			return nil
		},
	}
	deployRootCmd.Flags().String(deploymentManifestPath, defaultDeploymentManifestPath, "location of mazzaroth channel deployment manifest")
	return deployRootCmd
}
