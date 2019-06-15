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

## development

env2star is written in Go and has no dependencies other than the standard library.

```console
$ make
$ make test
```

## limitations

- You can't use `.`, `[`, or `]` in key names
- The top-level object can't be an array
