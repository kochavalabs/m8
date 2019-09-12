/**
 * This file encapsulates the basic reading and parsing of terminal input for
 * interactive contract-cli. It is fairly straight forward, and simply pulls
 * lines from stdin, parses them with a simple grammar, executes the correct
 * logic and prints the results to stdout.
*/
import nearley from 'nearley'
import grammar from './grammar/grammar.js'
import readline from 'readline'
import Debug from 'debug'

const debug = Debug('mazzaroth-cli:contract-io')

/**
 * Helper function for formatting an abiEntry and printing logging it with
 * console.log.
 *
 * @param abiEntry
 *
 * @return none
*/
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

/**
 * After constructing, you can call 'run' on a ContractIO object, which will
 * drop the user into an interactive CLI and only return when they exit the CLI.
 *
 * Wraps the readline library in some Mazzaroth specific logic to accomplish
 * this.
*/
class ContractIO {
  /**
   * The interactive contract CLI requires a contract client to work. This
   * constructor also creates an instance of the readline interface.
   *
   * @param contractClient The client used to make contract calls based on the
   *                       user input from readline.
  */
  constructor (contractClient) {
    debug('constructed contract IO with contractClient: %o', contractClient)
    this.rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
      prompt: 'Mazz> '
    })
    this.contractClient = contractClient
  }

  /**
   * Helper function that simply goes through all the functions in the abiJSON
   * provided by the contractClient and outputs them. This is called when the
   * user types 'abi' into the CLI to output what functions are available to
   * call.
  */
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

  /**
   * Executes a parsed contract function using the contract client, then logs
   * the result and prompts for the next line.
   *
   * @param functionName Name of the contract function to call.
   * @param args Arguments to the contract function in JSON format.
   *
   * @return none
  */
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

  /**
   * Called to initiate the interactive CLI for the user. Repeatedly calls
   * readline prompt parsing the results. First checks if the result is a single
   * command that exists on this class (for example abi). If so it executes that
   * function. Otherwise it is assumed that it is a properly parsed function
   * call then executed on the contract client.
  */
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
