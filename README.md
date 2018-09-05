# Terraform Provider

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Maintainers

This project is maintained by the developers at [Cedexis](https://www.cedexis.com/) (now part of Citrix&#174;).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.10.x
- [Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/cedexis/terraform-provider-citrixitm`

```sh
$ mkdir -p $GOPATH/src/github.com/cedexis; cd $GOPATH/src/github.com/cedexis
$ git clone git@github.com:cedexis/terraform-provider-citrixitm.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/cedexis/terraform-provider-citrixitm
$ make build
```

## Using the provider

**TODO**

## Contributing

Contributions are welcome. Please see [Contributing](./CONTRIBUTING.md). 
