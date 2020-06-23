# Mazzaroth-CLI

[![CircleCI](https://circleci.com/gh/kochavalabs/mazzaroth-cli.svg?style=svg)](https://circleci.com/gh/kochavalabs/mazzaroth-cli)

A node CLI tool that wraps the [mazzaroth-js](https://github.com/kochavalabs/mazzaroth-js)
node/contract clients and facilitates interaction with a Mazzaroth node. Gives
access to both the raw node rpc endpoints and abstracted access to the contract
on a given Node.

## Installation

The CLI can be installed by using npm.

```bash
npm install -g mazzaroth-cli
```

## Commands

Below is the list of commands provided by mazzaroth-cli. The more complex
commands will be explained in further detail below.

| Command | Description | Example |
| ------- | ----------- | ------- |
| transaction-call | Submits a call transaction to a mazzaroth node. Arguments are XDR formatted as  base64 encoded strings. | mazzaroth-cli transaction-call my_func -a 9uZz -a f1zsfABG7J |
| readonly-call | Submits a readonly call transaction to a mazzaroth node. Arguments are XDR formatted as  base64 encoded strings. | mazzaroth-cli readonly-call my_func -a 9uZz -a f1zsfABG7J |
| contract-update | Submits an update transaction to a mazzaroth node. The format of the argument is a path to a file containing contract wasm bytes. | mazzaroth-cli contract-update ./test/data/hello_world.wasm |
| permission-update | Submits a permission transaction to a mazzaroth node that allows another account to act on your behalf. The argument to this call is the public key of the account you would like to grant access to. | mazzaroth-cli permission-update 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c |
| transaction-lookup |  Looks up the current status and results of a transaction by ID. Argument is a transaction ID (256 bit hex value). | mazzaroth-cli transaction-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c |
| block-lookup | Looks up a Block using either a block ID as hex or block Number. | mazzaroth-cli block-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c |
| block-header-lookup | Looks up a Block Header using either a block ID as hex or block Number. | mazzaroth-cli block-header-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c |
| receipt-lookup | Looks up a transaction receipt. Argument is a transaction ID (256 bit hex value). | mazzaroth-cli receipt-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c |
| nonce-lookup | Looks up the current nonce for an account, Argument is an account ID (256 bit hex value). | mazzaroth-cli nonce-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c |
| account-lookup | Looks up the current information for an account, Argument is an account ID (256 bit hex value). | mazzaroth-cli account-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c |
| channel-lookup | Looks up the current information for a channel, Argument is what specifically to lookup about the channel. Current options: config/contract | mazzaroth-cli channel-lookup config |
| contract-cli | Drops into a contract cli where you can call contract functions interactively. | mazzaroth-cli contract-cli abi.json |
| [subscribe](https://mazzaroth.io/docs/4-Event_Subscription/3-Tools.md) | Subscribes to the receipts received by a readonly/standalone node. | mazzaroth-cli subscribe '{"receiptFilter": {}, "transactionFilter": {"configFilter":{}}}'|
| deploy | Helper for deploying a contract to a mazzaroth network. Takes a json config file, a sample config file can be found [here](https://github.com/kochavalabs/full-contract-example/blob/master/deploy.json) | mazzaroth-cli deploy ./deploy.json |
| xdr | Command used for converting between JSON and base64 representations of xdr objects. Also can be piped to from stdin. | mazzaroth-cli xdr Transaction '{"action": { "nonce": "3" } }' |


### Contract CLI

The call and lookup operations are relatively low level. The results need to be
interpreted from base64 strings or require multiple calls to complete. For
example to complete a 'transaction-call', you would need to look up an account
nonce, make the call, and finally lookup the results after execution. An example
of this being done (using  node.js) can be seen in the
[mazzaroth-js](https://github.com/kochavalabs/mazzaroth-js) repo.

This is cumbersome, so we've provided a further abstraction called the contract
client. This wraps the low level operations and gives the user access to their
contract through an rpc-like interface. We'll walk through how to drop into the
contract clients interactive CLI for a simple contract.

```bash
# The contract client requires the ABI json produced from our contract to run
# properly. Which will drop you into an interactive CLI. Normally this would
# be output by part of the mazzaroth build process, for this example we'll
# just pretend we have an abi json in the location it would normally be output
# to. This example can be found in our ore extensive example contract:
# https://github.com/kochavalabs/full-contract-example
mazzaroth-cli contract-cli contract/target/json/ExampleContract.json
Mazz>

# You can see the currently available functions by typing abi
Mazz> abi

Functions:
  args(string, string, string) -> uint32
  complex(Foo, Bar) -> string

ReadOnly Functions:
  simple() -> string

# And call them
Mazz> simple()
Hello World!
Mazz> args("one", "two", "three")
11
Mazz> complex('', '')
Error: Type not identified: Foo
    at nodeClient.nonceLookup.then.result
    at process._tickCallback
Mazz>
```

### Deploy

The deploy command is used to initialize a mazzaroth channel with its starting
configuration and contract. It also allows you to run some initial transactions
on startup if you desire.

```bash
# If you're just testing out a standalone node you can just echo an empty config
# into the deploy command to initialize a node to the default configuration.
echo '{}' | mazzaroth-cli deploy

# Or for a more complete configuration:
cat deploy.json
# {
#     "abi": {
#         "type": "file",
#         "value": "./contract/target/json/ExampleContract.json"
#     },
#     "channel-id": "0000000000000000000000000000000000000000000000000000000000000000",
#     "contract": "./contract/target/wasm32-unknown-unknown/release/contract.wasm",
#     "node-addr": "http://localhost:8080",
#     "on-behalf-of": "",
#     "owner": "3b6a27bcceb6a42d62a3a8d02a6f0d73653215771de243a63ac048a18b59da29",
#     "init-transactions": {
#         "initialize-configuration": [
#             {
#                 "args": [],
#                 "function_name": "setup",
#                 "sender": ""
#             }
#         ]
#     },
#     "xdr-types": "./xdrTypes.js"
# }
mazzaroth-cli deploy ./deploy.json
```

### XDR

Mazzaroth-cli also provides an xdr command that will allow easier conversion of
the common XDR types found in our [mazzaroth-xdr](https://github.com/kochavalabs/mazzaroth-xdr)
repository between base64 and the JSON representations. Some basic examples:

```bash

# Translate a basic empty transaction to base64 and back
mazzaroth-cli xdr Transaction '{}'
# AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA 

# You can also pipe to the xdr command from stdin
mazzaroth-cli xdr Transaction '{}' | mazzaroth-cli xdr Transaction --inputType base64
# {
#    "signature":"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
#    "signer":{
#       "enum":0,
#       "value":""
#    },
#    "action":{
#       "address":"0000000000000000000000000000000000000000000000000000000000000000",
#       "channelID":"0000000000000000000000000000000000000000000000000000000000000000",
#       "nonce":"0",
#       "category":{
#          "enum":0,
#          "value":""
#       }
#    }
# }

# Translate a ReceitSubscriptionResult from base64

mazzaroth-cli xdr ReceiptSubscriptionResult --inputType base64 AAAAAQwZ2DC+GfzdV7UPw15cE/P6KQawvWjsfMgpLh4i//8qAAAAAHfJJ+t0nVQb1/bhniZi1NoR7UNVy3X8ccPnhzSdwlPb
# {
#    "receipt":{
#       "status":1,
#       "stateRoot":"0c19d830be19fcdd57b50fc35e5c13f3fa2906b0bd68ec7cc8292e1e22ffff2a",
#       "result":""
#    },
#    "transactionID":"77c927eb749d541bd7f6e19e2662d4da11ed4355cb75fc71c3e787349dc253db"
# }
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
