# env2conf

[![Go Report Card](https://goreportcard.com/badge/github.com/andreykaipov/env2conf)](https://goreportcard.com/report/github.com/andreykaipov/env2conf)
[![CircleCI](https://img.shields.io/circleci/build/github/andreykaipov/env2conf/master.svg)](https://circleci.com/gh/andreykaipov/env2conf)
[![Wow Badges](https://img.shields.io/badge/wow-badges-blue.svg)](https://github.com/andreykaipov/env2conf)

env2conf converts environment variables into configuration files.

## usage

Available configuration outputs are:

- JSON (default)

  ```console
  $ env server.host=0.0.0.0 server.port=8080 env2conf -prefix server
  {
    "server": {
      "port": 8080,
      "host": "0.0.0.0"
    }
  }
  ```

- YAML
  ```console
  $ env fruits[0]=apple fruits[1]=banana fruits[2]=orange env2conf -prefix fruits -output yaml
  ---
  fruits:
    - banana
    - mango
    - pineapple
  ```

- TOML (inline)
  ```console
  $ env inputs.cpu[0]={} outputs.file[0].files[0]=stdout env2conf -prefix inputs,outputs -output toml
  outputs = {file = [{files = ["stdout"]}]}
  inputs = {cpu = [{}]}
  ```

## installation

Binaries are available for download from the [GitHub releases](https://github.com/andreykaipov/env2conf/releases) page.
For example:

```console
$ curl -Lo env2conf https://github.com/andreykaipov/env2conf/releases/download/v0.1.1/env2conf-0.1.1-linux-amd64
$ chmod +x env2conf
$ mv env2conf /usr/local/bin
```

Alternatively, if you have Go installed:

```console
$ go get github.com/andreykaipov/env2conf
```

## development

env2conf is written in Go and has no external dependencies. Make is the build tool:

```console
$ make
$ make test
```

## limitations

- You can't use `.`, `[`, or `]` in key names
- The top-level object can't be an array
