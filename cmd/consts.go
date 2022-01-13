package cmd

const (
	version                       = `0.0.1`
	cfgDir                        = `/.m8/`
	cfgName                       = `cfg.yaml`
	defaultBlockExpirationNumber  = 5
	defaultDeploymentManifestPath = `./m8/deployment.yaml`
	defaultTestManifestPath       = `./m8/test.yaml`
	// Flags/Env
	// Environment variables are expected to be ALL CAPS
	cfgPath                = `cfg-path`
	channelId              = `channel-id`
	address                = `address`
	transactionid          = `tx-id`
	headers                = `headers`
	header                 = `header`
	blockid                = `block-id`
	number                 = `number`
	height                 = `height`
	function               = `fn`
	args                   = `args`
	deploymentManifestPath = `deployment-manifest`
	testManifestPath       = `test-manifest`
	pause                  = `pause`
)
