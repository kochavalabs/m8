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

## License

[MIT](https://choosealicense.com/licenses/mit/)
