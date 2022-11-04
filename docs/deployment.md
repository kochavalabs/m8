# Deployment Configuration

Using m8, the command to deploy a contract is:

```Bash
m8 channel exec deployment --deployment-manifest deployment.yaml
```

Where `deployment.yaml` is the path to a manifest file containing the deployment
configuration for the contract. An example deployment.yaml can be found
[here](https://github.com/kochavalabs/m8/blob/develop/examples/deployment.yaml).

A basic deployment yaml should look like this:

```yaml
version: 0.0.1
type: deployment 
channel:
  version: 0.0.1
  id: 0000000000000000000000000000000000000000000000000000000000000000
  owner: 0000000000000000000000000000000000000000000000000000000000000000
  contract-file: "samplecontract.wasm"
  abi-file: "samplecontract.json"
gateway-node:
  address: https://localhost:6299 
deploy:
  name: sample-contract
  transactions:
    - tx: 
        function: "migration_database_1" 
        args: ["1","2","3","4"]
    - tx: 
        function: "migration_database_2" 
        args: ["1","2","3","4"]
    - tx: 
        function: "balanceof" 
        args: ["1","2","3","4"]
```

The first version is the version of the config file itself. There are two
types of config files that `m8` accepts (test and deployment) so for
this you should specify `deployment`.
The channel provides all the fields configuring your channel including:

- version - The version of the channel which should be incremented for new
deployments.
- id - The id of the channel which may be all 0s for testing purposes.
- owner - The public key id of the owner of the channel.
- contract-file - The path to the compiled Wasm contract file.
- abi-file - The path to the compiled json ABI.

The gateway-node address is the url of the Mazzaroth node to target for deployment,
which can be a locally running node.

The deploy section gives a name to the contract and can optionally be used to provide
a list of transactions to execute following the deployment.
