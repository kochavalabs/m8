package cmd

const (
	version                       = `0.0.1`
	channelIdLength               = 32
	pubKeyLength                  = 32
	privKeylength                 = 64
	cfgDir                        = `/.m8/`
	cfgName                       = `cfg.yaml`
	defaultDeploymentManifestPath = `./m8/deployment.yaml`
	defaultTestManifestPath       = `./m8/test.yaml`
	defaultChannelId              = `0000000000000000000000000000000000000000000000000000000000000000`
	defaultGatewayNodeAddress     = `http://localhost:6299`
	maxBlockExpirationRange       = 100

	// Flags/Env
	// Environment variables are expected to be ALL CAPS
	cfgPath            = `cfg-path`
	privateKey         = `private-key`
	publicKey          = `public-key`
	channelId          = `channel-id`
	channelAlias       = `channel-alias`
	channelAddress     = `channel-address`
	transactionId      = `tx-id`
	headers            = `headers`
	header             = `header`
	blockid            = `block-id`
	number             = `number`
	height             = `height`
	function           = `fn`
	arguments          = `args`
	deploymentManifest = `deployment-manifest`
	testManifest       = `test-manifest`
	pausechannel       = `pause`
)
