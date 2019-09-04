import nearley from 'nearley'
import grammar from './grammar/grammar.js'
import readline from 'readline'
import Debug from 'debug'

const debug = Debug('mazzaroth-cli:contract-io')

function outputAbiFunc (abiEntry) {
  let output = `  ${abiEntry.name}(`
  const types = abiEntry.inputs.map(x => x.type)
  output += types.join(', ')
  output += ')'
  if (abiEntry.outputs[0]) {
    output += ` -> ${abiEntry.outputs[0].type}`
  }
  console.log(output)
}

class ContractIO {
  constructor (contractClient) {
    debug('constructed contract IO with contractClient: %o', contractClient)
    this.rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
      prompt: 'Mazz> '
    })
    this.contractClient = contractClient
  }

  abi () {
    const functions = this.contractClient.abiJson.filter(x => x.type === 'function')
    const readonlys = this.contractClient.abiJson.filter(x => x.type === 'readonly')
    console.log()
    console.log('Functions: ')
    functions.forEach(outputAbiFunc)
    console.log()
    console.log('ReadOnly Functions: ')
    readonlys.forEach(outputAbiFunc)
    console.log()
    this.rl.prompt()
  }

  executeContractFunction (functionName, args) {
    if (!this.contractClient[functionName]) {
      throw new Error(`${functionName} is not a contract function`)
    }
    this.contractClient[functionName](...args).then(res => {
      console.log(res)
      this.rl.prompt()
    }).catch(e => {
      console.log(e)
      this.rl.prompt()
    })
  }

  run () {
    this.rl.prompt()

    this.rl.on('line', (line) => {
      try {
        const res = new nearley.Parser(nearley.Grammar.fromCompiled(grammar)).feed(line)
        if (res.results.length) {
          if (this[res.results[0]]) {
            this[res.results[0]]()
          } else {
            this.executeContractFunction(res.results[0].name, res.results[0].args)
          }
        } else {
          console.log(`Incomplete statement: "${line}"`)
          this.rl.prompt()
        }
      } catch (e) {
        console.log(e)
        this.rl.prompt()
      }
    }).on('close', () => {
      console.log('')
      console.log('peace bro~')
      process.exit(0)
    })
  }
}

export default ContractIO
