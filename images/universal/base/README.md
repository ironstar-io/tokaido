Base
======

Base is a common baseline for all Nozomi images and is
designed to offer the best possible developer experience while
ensuring high levels of security and sane defaults.

## Inclusions

Base includes the following components:

- zsh shell
- gcloud and aws command line interfaces
- TLS certificate bundles and OpenSSH client
- Network tools such as ping and telnet
- Text editors: Neovim and Nano
- Global user accounts and groups
- Basic Nozomi directory structure

As a general rule, if a component is needed by 3 or more
Nozomi images, then that component is included here.
