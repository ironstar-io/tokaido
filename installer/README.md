# 🚅 Tokaido Self-Installer

Tokaido is an automation platform for Drupal development environments on MacOS, Linux, and Windows

This installer is generally to be loaded onto a USB device or similar to facilitate an easy and fast install of Tokaido, composer cache and Docker images. Everything you need to get started with a new Tokaido environment!

At this stage internet connectivity is still required in order to download a few smaller files, however we are planning to make this installer fully offline in a future release.

## Requirements

Before you can use Tokaido, you need to have Docker installed on your systme.

### For MacOS

For MacOS, you need to install "Docker for Mac", which is in the included `Docker.dmg` file. 

### For Windows

Docker on Windows requires Windows Pro edition. Sorry, but Windows Home edition is not supported. 

- Install Docker Desktop by running the `Docker Desktop Installer.exe` included in this USB
- Install the Git Bash console by running `git-bash.exe`

## Installing Tokaido

In this directory, there should be three files corresponding to your OS type. 

- `tok-installer-macos`: Apple MacOS El Capitan and above
- `tok-installer-linux-amd64`: Linux distributions. Tested to this point in time with Ubuntu and Mint
- `tok-installer-windows-amd64`: Windows 10

After installing Docker, Executing the file that corresponds to your OS should open your default terminal and run the self-installer.

## Next Steps

You're now ready to go with Tokaido. Run `tok new` and start a new Drupal project!

- The Tokaido main README can be found here: https://github.com/ironstar-io/tokaido
- The full Tokaido documentation can be found here: https://docs.tokaido.io

## Talk to us!

If you have any feedback on Tokaido, or you're running into problems, we'd love
to hear from you.

For general support and queries, please visit us on the `#tokaido` channel in
the [offical Drupal Slack workspace](https://www.drupal.org/slack).

If you'd like to talk about commercial support arrangements for your team,
please [email us](tokaido@ironstar.io).

Found a bug or have a feature request? Please open a [GitHub Issue](https://github.com/ironstar-io/tokaido/issues/new/choose)

PRs are _more_ than welcome, but we suggest getting in touch if you'd like to
contribute any major feature. We have an open fortnightly planning session we
will be happy to invite you to.
