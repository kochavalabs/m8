# Test Configuration

Using m8, the command to test a contract is:

```Bash
m8 channel exec test --test-manifest test.yaml
```

Where `test.yaml` is the path to a manifest file containing the test
configuration for the contract. An example test.yaml can be found
[here](https://github.com/kochavalabs/m8/blob/develop/examples/test.yaml).

A basic test yaml should look like this:

```yaml
version: 0.0.1
type: test 
channel:
  version: 0.0.1
  id: 0000000000000000000000000000000000000000000000000000000000000000
  owner: 0000000000000000000000000000000000000000000000000000000000000000
  contract-file: "samplecontract.wasm"
  abi-file: "samplecontract.json"
gateway-node: 
  address: https://localhost:6299
tests :
  - name: test-foo
    reset: false 
    transactions:
      - tx:
        function: "foo"
        args: ["1","2","3","4"]
        receipt: 
          status: 1 
          result: "success"
      - tx:
        function: "bar"
        args: ["1","2","3","4"]
        receipt: 
          status: 1 
          result: "success"
```

The first version is the version of the config file itself. There are two
types of config files that `m8` accepts (test and deployment) so this
should be `test` for test config files.
The channel provides all the fields configuring your channel including:

- version - The version of the channel which should be incremented for new
deployments.
- id - The id of the channel which may be all 0s for testing purposes.
- owner - The public key id of the owner of the channel.
- contract-file - The path to the compiled Wasm contract file.
- abi-file - The path to the compiled json ABI.

The gateway-node address is the url of the Mazzaroth node to target for deployment,
which can be a locally running node.

The tests section gives a name to the test and can be used to provide
a list of transactions to execute following the deployment for testing.

The function and args are used to create the transaction and the result of
submitting the transaction is compared against the receipt values.
If the receipts do not match an error is reported.
