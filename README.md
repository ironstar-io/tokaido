# 🚅 Tokaido by Ironstar

[![CircleCI](https://circleci.com/gh/ironstar-io/tokaido.svg?style=shield)](https://circleci.com/gh/ironstar-io/tokaido)
[![GitHub license](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/ironstar-io/tokaido)

Tokaido is a Drupal environment launcher that _just works_. 

It creates Drupal environments in seconds and requires no manual configuration 
for your Drupal site to work out-of-the-box. 

[Click here to watch our full demo video](https://www.youtube.com/watch?v=pxktV9zQUhM&lc=z23nhfs54myvifnwn04t1aokg1km2r2d2ts4lrdilt4xrk0h00410) or just check out the latest quick-demo below.

[![Tokaido 1.3.0 Quickdemo](https://i.imgur.com/CLAW9If.png)](https://www.youtube.com/watch?v=nEb20jM31_8)

## Installation

Installation instructions are available on the [Tokaido website](https://tokaido.io/docs/getting-started/)

## Features

- Drupal 7 and 8 Support.
- Runs on MacOS, Linux, and Windows.
- Easy to use! Just type `tok up` for a zero-configuration environment.
- Fast! Most environments are ready in less than 30 seconds.
- Highly configurable with an easy to use config editor (`tok config`)
- Production-grade environments with Varnish and HAProxy included. 
- Easily add services like Mailhog, Solr, PHP Adminer, Xdebug, and more
- Commercial support available

## Why switch to Tokaido?

Tokaido was built out of frustration with how much time Drupal developers spend
trying to manage their local development environments. While other great tools
like DrupalVM and Lando have made it easier and easier for developers to build
repeatable Drupal environments, we felt there was still a lot of room for 
improvement. 

With Tokaido, we have shifted towards a more opinionated Drupal environment 
setup and coupled it with the same containers that we (Ironstar) run in 
enterprise Drupal hosting environments. 

So what makes Tokaido faster and more efficient?

- Using Docker instead of Vagrant enables faster, more light-weight environments
- We use Unison to sync files between your system and the Tokaido environment, so there's no performance hit from slow virtual disks
- The powerful `tok` CLI streamlines your set up and helps keep you out of config files
- Tokaido's inbuilt proxy enables verified SSL connections to your environment: https://local.tokaido.io:5154
- Built-in Varnish caching enables production-like testing of your code

## How does Tokaido compare?

|                                        | Tokaido                | Docker4Drupal   | Lando           |
|----------------------------------------|------------------------|-----------------|-----------------|
| Startup Time (excluding download)      | < 60 seconds           | < 60 seconds    | < 60 seconds    |
| Installation                           | `brew install tokaido` | Manual Download | Manual Download |
| Usability                              | Easy                   | Complex*        | Moderate+       |
| Works out-of-the-box**                 | Yes                    | No              | No              |
| Full Drupal/Drush SSH environment      | Yes                    | No              | No              |
| Production-ready containers            | Yes                    | No              | No              |
| Automated DB configuration             | Yes                    | No              | No              |
| Automated SSL configuration            | Yes                    | No              | No              |
| Modify PHP Runtime Config              | Yes                    | Yes             | Yes             |
| Run multiple environments easily       | Yes                    | No              | No              |
| Dev Tools - `yarn`, `npm`, `ruby`, etc | Yes                    | No              | No              |

Tokaido also ships with an incredibly powerful CLI that helps to eliminate the
need for you to manually manage config files. 

|                                  | Tokaido              | Docker4Drupal                         | Lando         |
|----------------------------------|----------------------|---------------------------------------|---------------|
| Powerful CLI                     | Yes                  | No                                    | Kinda         |
| Start new projects               | `tok new`            | -                                     | -             |
| Launch an environment            | `tok up`             | `docker-compose up -d`                | `lando start` |
| Edit configuration               | `tok config`         | -                                     | -             |
| Connect Drupal to database       | `tok up`             | -                                     | -             |
| Self-checks                      | `tok status`         | -                                     | -             |
| SSH into development environment | `ssh project.tok`    | -                                     | `lando ssh`++ |
| Run commands in an environment   | `tok exec "command"` | `docker-compose exec {container} ...` | -             |
| Reset Varnish cache              | `tok purge`          | -                                     | -             |
| Open site in browser             | `tok open`           | -                                     | -             |
| Open services in browser         | `tok open {service}` | -                                     | -             |
| Generate a Drupal hash satl      | `tok hash`           | -                                     | -             |

\* Docker4Drupal is controlled by a Docker Compose file and requires an indepth
understanding of Docker and Docker Compose in order to make changes. 

\+ Lando provides a helpful CLI that makes starting and managing environments
easier, but we still think Tokaido has it beat in this department. 

\+\+ Lando's 'ssh' environment is really just a docker-compose exec command, 
which is helpful for running commands but is not a full-featured dev environment

\*\* Nearly every Drupal project we've tested works with Tokaido without any 
special config. When testing Lando and Docker4Drupal, even the most basic Drupal
minimal installation required special config to get going. 

What about DrupalVM, DDev, and other tools? In the case of VM-based Drupal
environments, we feel these are all too slow for modern web development, with
some tools taking over 30 minutes to launch test environments. 

## What can't Tokaido do yet?

With Tokaido, we've tried to make it easy to get up and going with most Drupal
projects, but there are still some cases where Tokaido is improving:

- Multisite and Domain Access sites. These will work, but you'll need to manually set up your database. For domains, anything at *.local.tokaido.io will resolve to your local host. 
- PHP 5.6 and 7.2. Right now, Tokaido works with PHP 7.1 only. We have no plans to add support for PHP 5.6, but 7.2 support is coming soon.
- Xdebug works but is experimental. Please help us test it !

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
