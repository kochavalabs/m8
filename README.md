# Mazzaroth-CLI

A node CLI tool that wraps the [mazzaroth-js](https://github.com/kochavalabs/mazzaroth-js)
node/contract clients and facilitates interaction with a Mazzaroth node. Gives
access to both the raw node rpc endpoints and abstracted access to the contract
on a given Node.

## Installation

The CLI can be installed by using npm.

```bash
npm install -g mazzaroth-cli
```

## Usage

Execute 'mazzaroth-cli --help' for a list of commands. This CLI gives full
access to a Mazzaroth node through the commands defined by the rpc
Requests/Responses defined [here](https://github.com/kochavalabs/mazzaroth-xdr/blob/develop/idl/rpc.x).
For more context, an example of using the CLI with an actual contract can be
found [here](https://github.com/kochavalabs/full-contract-example).

## License

[MIT](https://choosealicense.com/licenses/mit/)
