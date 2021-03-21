# helmtrans
Yaml to helm 

## Usage:

### Yaml to Helm

```sh
➜ helmtrans  -h

helmtrans is a CLI library for Go that support yaml to helm.

Usage:
  helmtrans [command]

Available Commands:
  help        Help about any command
  version     Print the version number of helmtrans
  yamltohelm  Transform yaml to helm

Flags:
  -h, --help   help for helmtrans
```

use yamltohelm command

```sh
helmtrans yamltohelm -p [source path] -o [output path]
```
-o is optional param (default is output)


### Helm to Yaml

you can use [**schelm**](https://github.com/databus23/schelm) to render a helm manifest to a directory.

## Maintainer
- roancsu@163.com
- codeandcode0x@gmail.com