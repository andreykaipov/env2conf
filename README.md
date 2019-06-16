# env2star

env2star converts environment variables into configuration files.

## usage

Available configuration outputs are:

- JSON (default)

  ```console
  $ env server.host=0.0.0.0 server.port=8080 env2star -prefix server
  {
    "server": {
      "port": 8080,
      "host": "0.0.0.0"
    }
  }
  ```

- YAML
  ```console
  $ env fruits[0]=apple fruits[1]=banana fruits[2]=orange env2star -prefix fruits -output yaml
  ---
  fruits:
    - banana
    - mango
    - pineapple
  ```

- TOML (inline)
  ```console
  $ env inputs.cpu[0]={} outputs.file[0].files[0]=stdout env2star -prefix inputs,outputs -output toml
  outputs = {file = [{files = ["stdout"]}]}
  inputs = {cpu = [{}]}
  ```

## installation

Binaries are available for download from the [GitHub releases](https://github.com/andreykaipov/env2star/releases) page.
For example:

```console
$ curl -Lo env2star https://github.com/andreykaipov/env2star/releases/download/v0.1.0/env2star-0.1.0-linux-amd64
$ chmod +x env2star
$ mv env2star /usr/local/bin
```

Alternatively, if you have Go installed:

```console
$ go get github.com/andreykaipov/env2star
```

## development

env2star is written in Go and has no external dependencies. Make is the build tool:

```console
$ make
$ make test
```

## limitations

- You can't use `.`, `[`, or `]` in key names
- The top-level object can't be an array
