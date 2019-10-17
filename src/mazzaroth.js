/**
 * CLI tool that wraps the mazzaroth-js node and contract clients that
 * facilitates interaction with a Mazzaroth node. Gives access to both the raw
 * node rpc endpoints and abstracted access to the contract on a given Node.
*/
import path from 'path'
import { NodeClient, ContractClient } from 'mazzaroth-js'
import ContractIO from './contract-io.js'
import program from 'commander'
import fs from 'fs'
import { Schema } from 'mazzaroth-xdr'

const defaultChannel = '0'.repeat(64)

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
    options.host = options.host || 'http://localhost:8081'
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
  ]
]

// Used for transaction-call cli command, which specifically calls write
// functions on the contract. This requires function arguments.
const callArgs = []
const callOptions = [
  [
    '-a --args <args>',
    'Arguments to the transaction call as base64 encoded strings.',
    function (arg) {
      callArgs.push(arg)
    }
  ]
]

const transactionCallDesc = `
Submits a call transaction to a mazzaroth node. 
(https://github.com/kochavalabs/mazzaroth-xdr)

Examples:
  mazzaroth-cli transaction-call my_func -a 9uZz -a f1zsfABG7J
`
clientCommand('transaction-call', transactionCallDesc, transactionOptions.concat(callOptions),
  (val, options, client) => {
    const action = {
      channelID: options.channel_id || defaultChannel,
      nonce: (options.nonce || Math.floor(Math.random() * Math.floor(1000000000))).toString(),
      category: {
        enum: 1,
        value: {
          function: val,
          parameters: callArgs
        }
      }
    }
    client.transactionSubmit(action, options.on_behalf_of).then(res => {
      console.log(res.toJSON())
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

const readonlyCallDesc = `
Submits a readonly call transaction to a mazzaroth node. 
(https://github.com/kochavalabs/mazzaroth-xdr)

Examples:
  mazzaroth-cli readonly-call my_func -a 9uZz -a f1zsfABG7J
`
clientCommand('readonly-call', readonlyCallDesc, transactionOptions.concat(callOptions),
  (val, options, client) => {
    const call = {
      function: val,
      parameters: callArgs
    }
    client.readonlySubmit(call).then(res => {
      console.log(res.toJSON())
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

const updateOptions = [
  [
    '-s --schema <schema>',
    'File path to a binary representation of the contract schema.'
  ]
]

const contractUpdateDesc = `
Submits an update transaction to a mazzaroth node. The format of <val> is a path
to a file containing contract wasm bytes.
(https://github.com/kochavalabs/mazzaroth-xdr)

Examples:
  mazzaroth-cli contract-update ./test/data/hello_world.wasm --schema ./schema.xdr
`
clientCommand('contract-update', contractUpdateDesc, transactionOptions.concat(updateOptions),
  (val, options, client) => {
    fs.readFile(val, (err, data) => {
      let schemaObj = {
        tables: []
      }
      if (options.schema != null) {
        try {
          const schemaData = fs.readFileSync(options.schema)
          schemaObj = Schema().fromXDR(schemaData).toJSON()
        } catch (err) {
          console.log(err)
          return
        }
      }
      const action = {
        channelID: options.channel_id || defaultChannel,
        nonce: (options.nonce || Math.floor(Math.random() * Math.floor(1000000000))).toString(),
        category: {
          enum: 2,
          value: {
            contract: data.toString('base64'),
            schema: schemaObj
          }
        }
      }
      if (err) throw err
      client.transactionSubmit(action, options.on_behalf_of).then(res => {
        console.log(res.toJSON())
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
    const permType = Number(options.perm_type) || 0
    const action = {
      channelID: options.channel_id || defaultChannel,
      nonce: (options.nonce || Math.floor(Math.random() * Math.floor(1000000000))).toString(),
      category: {
        enum: 3,
        value: {
          key: val,
          action: permType
        }
      }
    }
    client.transactionSubmit(action, options.on_behalf_of).then(res => {
      console.log(res.toJSON())
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
      console.log(res.toJSON())
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
        console.log(res.toJSON())
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
      console.log(res.toJSON())
    })
      .catch(error => {
        if (error.response) {
          console.log(error.response.data)
        } else {
          console.log(error)
        }
      })
  })

const nonceLookupDesc = `
Looks up the current nonce for an account, Val is an account ID (256 bit hex value).

Examples:
  mazzaroth-cli nonce-lookup 3a547668e859fb7b112a1e2dd7efcb739176ab8cfd1d9f224847fce362ebd99c
`
clientCommand('nonce-lookup', nonceLookupDesc, [],
  (val, options, client) => {
    client.publicKey = Buffer.from(val, 'hex')
    client.nonceLookup().then(res => {
      console.log(res.toJSON())
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
      console.log(res.toJSON())
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

program.on('command:*', function (command) {
  program.help()
})

program.parse(process.argv)
