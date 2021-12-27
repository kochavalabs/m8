package cmd

const (
	defaultCfgPath                = `$HOME/.m8`
	defaultDeploymentManifestPath = `./m8/deployment.yaml`
	defaultTestManifestPath       = `.m8/test.yaml`
	// Flags/Env
	// Environment variables are expected to be ALL CAPS
	cfgPath                = `cfg`
	channelId              = `channel-id`
	address                = `address`
	transactionid          = `tx-id`
	headers                = `headers`
	blocks                 = `blocks`
	blockid                = `block-id`
	number                 = `number`
	height                 = `height`
	deploymentManifestPath = `deployment-manifest`
	testManifestPath       = `test-manifest`
)
