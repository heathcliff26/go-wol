[![CI](https://github.com/heathcliff26/go-wol/actions/workflows/ci.yaml/badge.svg?event=push)](https://github.com/heathcliff26/go-wol/actions/workflows/ci.yaml)
[![Coverage Status](https://coveralls.io/repos/github/heathcliff26/go-wol/badge.svg)](https://coveralls.io/github/heathcliff26/go-wol)
[![Editorconfig Check](https://github.com/heathcliff26/go-wol/actions/workflows/editorconfig-check.yaml/badge.svg?event=push)](https://github.com/heathcliff26/go-wol/actions/workflows/editorconfig-check.yaml)
[![Generate go test cover report](https://github.com/heathcliff26/go-wol/actions/workflows/go-testcover-report.yaml/badge.svg)](https://github.com/heathcliff26/go-wol/actions/workflows/go-testcover-report.yaml)
[![Renovate](https://github.com/heathcliff26/go-wol/actions/workflows/renovate.yaml/badge.svg)](https://github.com/heathcliff26/go-wol/actions/workflows/renovate.yaml)

# go-wol

This is a simple utility for sending Wake-On-Lan magic packet to clients in the local network.
It can be used directly via the cli, or remotely via a web interface.

## Table of Contents

- [go-wol](#go-wol)
  - [Table of Contents](#table-of-contents)
  - [Usage](#usage)
    - [CLI Args](#cli-args)
    - [Using the image](#using-the-image)
    - [Image location](#image-location)
    - [Tags](#tags)
  - [Credit](#credit)

## Usage

### CLI Args
```
$ go-wol help
go-wol power on other devices on the network via Wake-on-Lan

Usage:
  go-wol [flags]
  go-wol [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  server      Serve a frontend via gui
  version     Print version information and exit
  wol         Send a magic packet to the given mac address

Flags:
  -h, --help   help for go-wol

Use "go-wol [command] --help" for more information about a command.
```

### Using the image

When using the container image, please note that the server needs to run with `--net host` to send the magic packets.
```
$ podman run -d -net host -v /path/to/config.yaml:/config/config.yaml ghcr.io/heathcliff26/go-wol:latest
```

### Image location

| Container Registry                                                                                     | Image                                      |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------ |
| [Github Container](https://github.com/users/heathcliff26/packages/container/package/go-wol) | `ghcr.io/heathcliff26/go-wol`   |
| [Docker Hub](https://hub.docker.com/repository/docker/heathcliff26/go-wol)                  | `docker.io/heathcliff26/go-wol` |

### Tags

There are different flavors of the image:

| Tag(s)      | Description                                                 |
| ----------- | ----------------------------------------------------------- |
| **latest**  | Last released version of the image                          |
| **rolling** | Rolling update of the image, always build from main branch. |
| **vX.Y.Z**  | Released version of the image                               |

## Credit

The css is the free to use template from w3schools: [Kitchen Sink/W3.CSS Demo Template](https://www.w3schools.com/w3css/tryw3css_templates_black.htm)
