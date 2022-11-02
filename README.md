# m8

The `m8` (Pronounced m-ate) command line tool lets you interact with Mazzaroth nodes.
For configuration, `m8` looks for a file named `cfg.yaml` in the `$HOME/.m8/` directory.
However, you can specify other setting a env variable of `M8_CFG_PATH` or setting the global `--cfg-path` flag.

## Syntax

m8 follows the following syntax to run commands from your terminal:

```Bash
m8 [resource] [verb] [noun] [flags]
```

Where `resource` includes:

- cfg - Commands involving the configuration of m8.
- channel - Commands involving the interaction of a Mazzaroth channel.

If you need help, run `m8 help` from your terminal.

## Testing and Deployment

Two major functions of m8 include deploying a contract to a Mazzaroth channel and running
tests against a Mazzaroth channel.
There are special configuration files that can be created to use these commands.
For details on either of the configuration manifests see the documentation in the `/docs` directory.
