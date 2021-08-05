/**
 * CLI tool that wraps the mazzaroth-js node and contract clients that
 * facilitates interaction with a Mazzaroth node. Gives access to both the raw
 * node rpc endpoints and abstracted access to the contract on a given Node.
*/
import path from 'path'
import { NodeClient, ContractClient, ReceiptSubscribe, XDRtoJSON, JSONtoXDR } from 'mazzaroth-js'
import ContractIO from './contract-io.js'
import program from 'commander'
import fs from 'fs'
import { sha3_256 as sha3256 } from 'js-sha3'
require('app-module-path').addPath(path.resolve(`${__dirname}/../node_modules`))

const defaultChannel = '0'.repeat(64)
const defaultAddr = 'http://localhost:8081'
const defaultOwner = '3b6a27bcceb6a42d62a3a8d02a6f0d73653215771de243a63ac048a18b59da29'
const defaultSender = '0'.repeat(64)
const defaultExpiration = 5
const defaultVersion = '0.1'

/**
 * Many of the node client commands have similar options. This just wraps the
 * the common logic for these commands.
 *
 * @param command Name of the command in the cli i.e. 'nonce-lookup'
 * @param desc Description of the command
 * @param opts Any additional options beyond host and priv_key
 * @param action Function that is the actual logic for the specific command,
 *               accepts the '<val>', options, and constructed client as args.
 *
 * @return none
*/
const clientCommand = (command, desc, opts, action) => {
  let cmd = program.command(`${command} <val>`)
  cmd.description(desc)
    .option('-h --host <s>',
      'Web address of the host node default: "http://localhost:8081"')
    .option('-k --priv_key <s>',
      'Priv key hex to sign contract with.')

  for (const opt of opts) {
    if (opt[2]) {
      cmd.option(opt[0], opt[1], opt[2])
    } else {
      cmd.option(opt[0], opt[1])
    }
  }
  cmd.action(function (val, options) {
    options.host = options.host || defaultAddr
    const client = new NodeClient(options.host, options.priv_key)
    action(val, options, client)
  })
}

// Additional options for a 'transaction' type command. Used by any non-readonly
// transactions.
const transactionOptions = [
  [
    '-c --channel_id <s>',
    'Base64 channel ID to send transaction to.'
  ],
  [
    '-n --nonce <s>',
    'Nonce to sent the request with.'
  ],
  [
    '-o --on_behalf_of <s>',
    'Account to send the transaction as.'
  ],
  [
    '-e --transaction_expire_after <s>',
    'transaction expires if not included in the next <s> blocks.'
  ]
]

// Used for transaction-call cli command, which specifically calls write
// functions on the contract. This requires function arguments.
const callArgs = []
const callOptions = [
  [
    '-a --args <args>',
    'Arguments to the transaction call are strings or json for complex structs.',
    function (arg) {
      callArgs.push(arg)
    }
  ]
]

const transactionCallDesc = `
Submits a call transaction to a mazzaroth node. 
(https://github.com/kochavalabs/mazzaroth-xdr)

Examples:
  mazzaroth-cli transaction-call my_func -a 'arg_one' -a 'arg_two'
`
clientCommand('transaction-call', transactionCallDesc, transactionOptions.concat(callOptions),
  (val, options, client) => {
    const blockExpiration = defaultExpiration | options.transaction_expire_after
    const action = {
      channelID: options.channel_id || defaultChannel,
      nonce: (options.nonce || Math.floor(Math.random() * Math.floor(1000000000))).toString(),
      category: {
        enum: 1,
        value: {
          function: val,
          arguments: callArgs
        }
      }
    }
    client.blockHeightLookup().then(res => {
      action.blockExpirationNumber = (parseInt(blockExpiration) + parseInt(res)).toString()
      client.transactionSubmit(action, options.on_behalf_of).then(res => {
        console.log(JSON.stringify(res.toJSON()))
      })
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

// Version of the contract being deployed
const conOptions = [
  [
    '-v --contract_version <args>',
    'version number for the contract'
  ],
  [
    '-a --abi <args>',
    'ABI JSON file to include with contract. Required.'
  ],
  [
    '-e --transaction_expire_after <s>',
    'transaction expires if not included in the next <s> blocks.'
  ]
]

const contractUpdateDesc = `
Submits an update transaction to a mazzaroth node. The format of <val> is a path
to a file containing contract wasm bytes.
(https://github.com/kochavalabs/mazzaroth-xdr)

Examples:
  mazzaroth-cli contract-update ./test/data/hello_world.wasm
`
clientCommand('contract-update', contractUpdateDesc, transactionOptions.concat(conOptions),
  (val, options, client) => {
    if (options.abi) {
      fs.readFile(options.abi, (err, data) => {
        if (err) {
          console.log('could not read file: ' + val)
          return
        }
        const abiJSON = JSON.parse(data.toString('ascii'))

        const blockExpiration = defaultExpiration | options.transaction_expire_after
        fs.readFile(val, (err, data) => {
          const action = {
            channelID: options.channel_id || defaultChannel,
            nonce: (options.nonce || Math.floor(Math.random() * Math.floor(1000000000))).toString(),
            blockExpirationNumber: blockExpiration.toString(),
            category: {
              enum: 2,
              value: {
                enum: 1,
                value: {
                  contractBytes: data.toString('base64'),
                  abi: abiJSON,
                  contractHash: sha3256.create().update(data.buffer).hex(),
                  version: options.contract_version || '0.1.0'
                }
              }
            }
          }

          if (err) throw err

          client.blockHeightLookup().then(res => {
            action.blockExpirationNumber = (parseInt(blockExpiration) + parseInt(res)).toString()
            console.log(action)
            client.transactionSubmit(action, options.on_behalf_of).then(res => {
              console.log(JSON.stringify(res.toJSON()))
            })
          })
            .catch(error => {
              if (error.response) {
                console.log(error.response.data)
              } else {
                console.log(error)
              }
            })
        })
      })
    } else {
      console.log('Please provide an ABI JSON file')
    }
  })

// Command option specific to the granting/revoking of permission.
const permOptions = [
  [
    '-p --perm_type <args>',
    'Permission type. 0: revoke, 1: grant'
  ]
]

const permissionUpdateDesc = `
Submits a permission transaction to a mazzaroth node that allows another account
to act on your behalf. The argument to this call is the public key of the
account you would like to grant access to.

Examples:
  mazzaroth-cli permission-update 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c
`
clientCommand('permission-update', permissionUpdateDesc, transactionOptions.concat(permOptions),
  (val, options, client) => {
    const blockExpiration = defaultExpiration | options.transaction_expire_after
    const permType = Number(options.perm_type) || 0
    const action = {
      channelID: options.channel_id || defaultChannel,
      nonce: (options.nonce || Math.floor(Math.random() * Math.floor(1000000000))).toString(),
      blockExpirationNumber: blockExpiration.toString(),
      category: {
        enum: 2,
        value: {
          enum: 3,
          value: {
            key: val,
            action: permType
          }
        }
      }
    }
    client.blockHeightLookup().then(res => {
      action.blockExpirationNumber = (parseInt(blockExpiration) + parseInt(res)).toString()
      client.transactionSubmit(action, options.on_behalf_of).then(res => {
        console.log(JSON.stringify(res.toJSON()))
      })
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

const transactionLookupDesc = `
Looks up the current status and results of a transaction by ID. Val is  a
transaction ID (256 bit hex value).

Examples:
  mazzaroth-cli transaction-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c
`
clientCommand('transaction-lookup', transactionLookupDesc, [],
  (val, options, client) => {
    client.transactionLookup(val).then(res => {
      console.log(JSON.stringify(res.toJSON()))
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

/**
 * block-lookup and block-header-lookup are so similar that I made a function
 * that handles both.
 *
 * @param lookupFunc Function name on the node client
 * @param cmd CLI command for the function.
 * @param desc We need 'Block' or 'Block Header' in the description.
 *
 * @return none
*/
function blockLookupCommand (lookupFunc, cmd, desc) {
  const blockLookupDesc = `
Looks up a ${desc} using either a block ID as hex or block Number.
Examples:
  mazzaroth-cli ${cmd} 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c
  mazzaroth-cli ${cmd} 1000
  `
  clientCommand(cmd, blockLookupDesc, [],
    (val, options, client) => {
      const possibleInt = parseInt(val)
      if (!isNaN(possibleInt) && possibleInt.toString() === val) {
        val = possibleInt
      }
      client[lookupFunc](val).then(res => {
        console.log(JSON.stringify(res.toJSON()))
      })
        .catch(error => {
          if (error.response) {
            console.log(error.response.data)
          } else {
            console.log(error)
          }
        })
    })
}
blockLookupCommand('blockLookup', 'block-lookup', 'Block')
blockLookupCommand('blockHeaderLookup', 'block-header-lookup', 'Block Header')

const receiptLookupDesc = `
Looks up a transaction receipt. Val is a transaction ID (256 bit hex value).

Examples:
  mazzaroth-cli receipt-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c
`
clientCommand('receipt-lookup', receiptLookupDesc, [],
  (val, options, client) => {
    client.receiptLookup(val).then(res => {
      console.log(JSON.stringify(res.toJSON()))
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

const accountLookupDesc = `
Looks up the current information for an account, Val is an account ID (256 bit hex value).

Examples:
  mazzaroth-cli account-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c
`
clientCommand('account-lookup', accountLookupDesc, [],
  (val, options, client) => {
    client.publicKey = Buffer.from(val, 'hex')
    client.accountInfoLookup().then(res => {
      console.log(JSON.stringify(res.toJSON()))
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

const channelLookupDesc = `
Looks up the current information for a channel, Val is what specifically to
lookup about the channel. Current options:

'config': ContractChannelConfig
'contract': Contract (bytes and version)


Examples:
  mazzaroth-cli channel-lookup config
`
clientCommand('channel-lookup', channelLookupDesc, [],
  (val, options, client) => {
    const valLookup = { 'contract': 1, 'config': 2 }
    client.channelInfoLookup(valLookup[val]).then(res => {
      console.log(JSON.stringify(res.toJSON()))
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

const cliOptions = [
  [
    '-x --xdr_types <s>',
    'Custom struct types javascript file (made with xdrgen)'
  ],
  [
    '-c --channel_id <s>',
    'Base64 channel ID to send transaction to.'
  ],
  [
    '-o --on_behalf_of <s>',
    'Account to send the transaction as.'
  ]
]

const contractCliDesc = `
Drops into a contract cli where you can call contract functions interactively.

Examples:
  mazzaroth-cli contract-cli abi.json
`
clientCommand('contract-cli', contractCliDesc, cliOptions,
  (val, options, client) => {
    fs.readFile(val, (err, data) => {
      let xdrTypes = {}
      if (options.xdr_types) {
        xdrTypes = require(path.resolve(options.xdr_types))
      }
      if (err) {
        console.log('could not read file: ' + val)
        return
      }
      const abiJSON = JSON.parse(data.toString('ascii'))
      const contractClient = new ContractClient(abiJSON, client, xdrTypes, options.channel_id, options.on_behalf_of)
      const io = new ContractIO(contractClient)
      io.run()
    })
  })

const subCmd = program.command('subscribe [val]')
const subCmdDescription = `
Subscribes to the receipts received by a readonly/standalone node.

Examples:
  mazzaroth-cli subscribe '{"receiptFilter": {}, "transactionFilter": {"configFilter":{}}}'
`
subCmd.description(subCmdDescription).option('-h --host <s>', 'Web address of the host node default: "localhost:8081"')
subCmd.action(function (val, options) {
  options.host = options.host || 'localhost:8081'
  val = val || '{}'
  ReceiptSubscribe(options.host, JSON.parse(val), (result) => { console.log(result) })
  process.stdin.setRawMode(true)
  process.stdin.resume()
  process.stdin.on('data', process.exit.bind(process, 0))
})

const deployCmd = program.command('deploy [input]')
const deployCmdDescription = `
Helper for deploying a contract to a mazzaroth network. Takes a json config file,
a sample config file can be found at
https://github.com/kochavalabs/full-contract-example/blob/master/deploy.json

Examples:
  mazzaroth-cli deploy ./deploy.json
  echo '{}' | mazzaroth-cli deploy
`

deployCmd.description(deployCmdDescription)
  .option('-h --host <s>',
    'Web address of the host node default: "http://localhost:8081"')
  .option('-t --timeout <number>', 'Timeout before transaction is considered failed in ms default: 3000')
deployCmd.action(async function (input, options) {
  let config = {}
  if (stdin) {
    config = JSON.parse(stdin)
  } else {
    config = JSON.parse(fs.readFileSync(input))
  }
  const channel = config['channel-id'] || defaultChannel
  const version = config['contract-version'] || defaultVersion
  const owner = config['owner'] || defaultOwner
  const txExpiration = config['transaction-expire-after'] | defaultExpiration
  let host = options.host || config['host']
  host = host || defaultAddr

  const sender = config['sender'] || defaultSender
  const client = new NodeClient(host, sender)
  const blockHeight = await client.blockHeightLookup()
  const blockExpirationNumber = blockHeight + txExpiration

  const configAction = {
    channelID: channel,
    nonce: '3',
    blockExpirationNumber: blockExpirationNumber.toString(),
    category: {
      enum: 2,
      value: {
        enum: 2,
        value: {
          owner: owner,
          admins: []
        }
      }
    }
  }

  const timeout = options.timeout || 3000
  const configRes = await client.transactionForReceipt(configAction, null, timeout).then(x => x.toJSON())
  console.log(configRes)
  // If they didn't set an initial contract, exit after the config action.
  if (config['contract'] === undefined) {
    return
  }

  let wasmFile
  const contractRootDir = path.dirname(input)
  if (path.isAbsolute(config['contract'])) {
    wasmFile = fs.readFileSync(config['contract'])
  } else {
    const wasmPath = path.join(path.dirname(input), config['contract'])
    wasmFile = fs.readFileSync(wasmPath)
  }
  const abiConf = config['abi']
  let abi = abiConf['value']
  if (abiConf['type'] === 'file') {
    if (path.isAbsolute(abiConf['value'])) {
      abi = JSON.parse(fs.readFileSync(abiConf['value']))
    } else {
      const abiPath = path.join(contractRootDir, abiConf['value'])
      abi = JSON.parse(fs.readFileSync(abiPath))
    }
  }
  const blockExpirationNumber2 = blockExpirationNumber + 1
  const action = {
    channelID: channel,
    nonce: '10',
    blockExpirationNumber: blockExpirationNumber2.toString(),
    category: {
      enum: 2,
      value: {
        enum: 1,
        value: {
          contractBytes: wasmFile.toString('base64'),
          abi: abi,
          contractHash: sha3256.create().update(wasmFile.buffer).hex(),
          version: version
        }
      }
    }
  }

  const contractRes = await client.transactionForReceipt(action, null, timeout).then(x => x.toJSON())
  console.log(contractRes)
  let xdrTypes = {}
  if (config['xdr-types']) {
    if (path.isAbsolute(config['xdr-types'])) {
      xdrTypes = require(path.resolve(config['xdr-types']))
    } else {
      const xdrTypesPath = path.join(contractRootDir, config['xdr-types'])
      xdrTypes = require(path.resolve(xdrTypesPath))
    }
  }
  const transactions = config['init-transactions']
  for (const txName in transactions) {
    const txSet = config['init-transactions'][txName]
    for (const txIndex in txSet) {
      const tx = txSet[txIndex]
      const sender = tx['sender'] || defaultSender
      const client = new NodeClient(host, sender)
      const contractClient = new ContractClient(abi, client, xdrTypes, channel, null, timeout)
      const functionName = tx['function_name']
      const result = await contractClient[functionName](txExpiration, ...tx['args'].map(x => {
        if (typeof x === 'object' && x !== null) {
          return JSON.stringify(x)
        }
        return x
      }))
      console.log(`Transaction run: ${functionName}`)
      console.log(`Result: ${JSON.stringify(result)}`)
    }
  }
})

var stdin = ''
const xdrCmd = program.command('xdr <type> [input]')
const xdrCmdDescription = `
Command used for converting between JSON and base64 representations of xdr
objects. Also can be piped to from stdin.

Examples:
  mazzaroth-cli xdr Transaction '{"action": { "nonce": "3" } }'
  echo '{"action": { "nonce": "3" } }' | mazzaroth-cli xdr Transaction
`

xdrCmd.description(xdrCmdDescription)
  .option('-i --inputType <s>',
    'Input type to convert from, defaults to JSON other option is base64')
xdrCmd.action(async function (type, input, options) {
  if (stdin) {
    input = stdin
  }
  if (options.inputType === 'base64') {
    console.log(XDRtoJSON(input, type))
  } else {
    console.log(JSONtoXDR(input, type))
  }
})

program.on('command:*', function (command) {
  program.help()
})

if (process.stdin.isTTY) {
  program.parse(process.argv)
} else {
  process.stdin.on('readable', function () {
    var chunk = this.read()
    if (chunk !== null) {
      stdin += chunk
    }
  })
  process.stdin.on('end', function () {
    program.parse(process.argv)
  })
}
