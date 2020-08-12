## inlets is a Cloud Native Tunnel written in Go

<img src="docs/inlets-logo-sm.png" width="150px">

Expose your local endpoints to the Internet or to another network, traversing firewalls, proxies, and NAT.


[![Build Status](https://travis-ci.com/inlets/inlets.svg?branch=master)](https://travis-ci.com/inlets/inlets)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/inlets/inlets)](https://goreportcard.com/report/github.com/inlets/inlets)
[![Documentation](https://godoc.org/github.com/inlets/inlets?status.svg)](http://godoc.org/github.com/inlets/inlets)
![GitHub All Releases](https://img.shields.io/github/downloads/inlets/inlets/total)

[Follow @inletsdev on Twitter](https://twitter.com/inletsdev)

[English](./README.md) | [中文文档](./README_CN.md)

## Intro

inlets&reg; combines a reverse proxy and websocket tunnels to expose your internal and development endpoints to the public Internet via an exit-node. An exit-node may be a 5-10 USD VPS or any other computer with an IPv4 IP address. You can also tunnel services without exposing them on the Internet, making inlets a suitable replacement for a VPN.

Why do we need this project? Similar tools such as [ngrok](https://ngrok.com/) or [Argo Tunnel](https://developers.cloudflare.com/argo-tunnel/) from [Cloudflare](https://www.cloudflare.com/) are closed-source, have limits built-in, can work out expensive, and have limited support for arm/arm64. Ngrok is also often banned by corporate firewall policies meaning it can be unusable. Other open-source tunnel tools are designed to only set up a single static tunnel. inlets aims to dynamically bind and discover your local services to DNS entries with automated TLS certificates to a public IP address over a websocket tunnel.

When combined with SSL - inlets can be used with any corporate HTTP proxy which supports `CONNECT`.

![Conceptual diagram](docs/inlets.png)

*Conceptual diagram for inlets*

### Who is behind this project?

inlets is brought to you by [Alex Ellis](https://twitter.com/alexellisuk). Alex is a [CNCF Ambassador](https://www.cncf.io/people/ambassadors/) and the founder of [OpenFaaS](https://github.com/openfaas/faas/).

> [OpenFaaS&reg;](https://github.com/openfaas/faas) makes it easy for developers to deploy event-driven functions and microservices to Kubernetes without repetitive, boiler-plate coding. Package your code or an existing binary in a Docker image to get a highly scalable endpoint with auto-scaling and metrics. The project has around 21k GitHub stars, over 270 contributors and a growing number of end-users in production.

#### SWAG and merchandise

Head over to the new [OpenFaaS Ltd SWAG store](https://store.openfaas.com/) to get your very own branded inlets hoodie, t-shirt and mug.

<img src="https://pbs.twimg.com/media/EQuxmEJWoAAP0Ga?format=jpg&name=small" width=300>

#### Insiders subscription

Buy an [Insiders Subscription](https://insiders.openfaas.io/) to get regular updates on the inlets, openfaas, k3sup, arkade along with hints and tips on Go, Docker and Kubernetes.

## License & terms

**Important**

Developers wishing to use inlets within a corporate network are advised to seek approval from their administrators or management before using the tool. By downloading, using, or distributing inlets, you agree to the [LICENSE](./LICENSE) terms & conditions. No warranty or liability is provided.

### Backlog & goals

#### Completed

* automatically create endpoints on exit-node based upon client definitions
  * multiplex sites on same port and websocket through the use of DNS / host entries
* link encryption using SSL over websockets (`wss://`)
* authentication for server and client
* automatic reconnect
* native multi-arch with ARMHF/ARM64 support
* Dockerfile and Kubernetes YAML files
* discover and implement `Service` type of `LoadBalancer` for Kubernetes - [inlets-operator](https://github.com/inlets/inlets-operator)
* tunnelling websocket traffic in addition to HTTP(s)
* [get a logo for the project](https://github.com/inlets/inlets/issues/46)
* TLS certificates when used with a reverse proxy like Nginx, Caddy or Traefik

#### Inlets PRO

The following features / use-cases are covered by [inlets.pro](https://github.com/inlets/inlets-pro).

* Tunnel L4 TCP traffic in addition to HTTP/s at L7
* Automated TLS - including via inletsctl/inlets-operator
* Commercial services & support

### Status

Unlike HTTP 1.1 which follows a synchronous request/response model websockets use an asynchronous pub/sub model for sending and receiving messages. This presents a challenge for tunneling a *synchronous protocol* over an *asynchronous bus*.

inlets 2.0 introduces performance enhancements and leverages parts of the Kubernetes and Rancher API. It uses the same tunnelling packages that enable node-to-node communication in [Rancher's k3s project](https://k3s.io). It is suitable for development and may be useful in production. Before deploying `inlets` into production, it is advised that you do adequate testing.

Feel free to open issues if you have comments, suggestions or contributions.

* The tunnel link is secured via `--token` flag using a shared secret
* The default configuration uses websockets without SSL `ws://`, but to enable encryption you could enable SSL `wss://`
* A timeout for requests can be configured via args on the server
* ~~The upstream URL has to be configured on both server and client until a discovery or service advertisement mechanism is added~~ The client can advertise upstream URLs, which it can serve
* The tunnel transport is wrapped by default which strips CORS headers from responses, but you can disable it with the `--disable-transport-wrapping` flag on the server

### inlets projects

Inlets is a Cloud Native Tunnel and is [listed on the Cloud Native Landscape](https://landscape.cncf.io/category=service-proxy&format=card-mode&grouping=category&sort=stars) under *Service Proxies*.

* [inlets](https://github.com/inlets/inlets) - Cloud Native Tunnel for L7 / HTTP traffic written in Go
* [inlets-pro](https://github.com/inlets/inlets-pro-pkg) - Cloud Native Tunnel for L4 TCP
* [inlets-operator](https://github.com/inlets/inlets-operator) - Public IPs for your private Kubernetes Services and CRD
* [inletsctl](https://github.com/inlets/inletsctl) - Automate the cloud for fast HTTP (L7) and TCP (L4) tunnels

## Get inlets

You can install the CLI with a `curl` utility script, `brew` or by downloading the binary from the releases page. Once installed you'll get the `inlets` command.

### Install the CLI

> Note: `inlets` is made available free-of-charge, but you can support its ongoing development and sign up for updates through [GitHub Sponsors](https://github.com/sponsors/alexellis/) 💪

Utility script with `curl`:

```bash
# Install to local directory
curl -sLS https://get.inlets.dev | sh

# Install to /usr/local/bin/
curl -sLS https://get.inlets.dev | sudo sh
```

Via `brew`:

```bash
brew install inlets
```

> Note: the `brew` distribution is maintained by the brew team, so it may lag a little behind the GitHub release.

Binaries are made available on the [releases page](https://github.com/inlets/inlets/releases) for Linux (x86_64, armhf & arm64), Windows (experimental), and for Darwin (MacOS). You will also find SHA checksums available if you want to verify your download.

Windows users are encouraged to use [git bash](https://git-scm.com/downloads) to install the inlets binary.

## Using inlets

### Video demo

Using inlets I was able to set up a public endpoint (with a custom domain name) for my JavaScript & Webpack [Create React App](https://github.com/facebook/create-react-app).

[![https://img.youtube.com/vi/jrAqqe8N3q4/hqdefault.jpg](https://img.youtube.com/vi/jrAqqe8N3q4/maxresdefault.jpg)](https://youtu.be/jrAqqe8N3q4)

### Quickstart tutorial

You can run inlets between any two computers with connectivity, these could be containers, VMs, bare metal or even "loop-back" on your own laptop.

Try the [quickstart tutorial now](./docs/quickstart.md) on your local computer.

### Documentation & tutorials

inlets and inlets PRO now has a dedicated documentation site:

Official docs: [docs.inlets.dev](https://docs.inlets.dev)

Other inlets OSS documentation & tutorials:

* Tutorial: [HTTPS for your local endpoints with inlets and Caddy](https://blog.alexellis.io/https-inlets-local-endpoints/)
* Docs: [Quickstart tutorial on your laptop](./docs/quickstart.md)
* Docs: [Inlets & Kubernetes recipes](./docs/kubernetes.md)
* Docs: [Run Inlets on a VPS](./docs/vps.md)
* Tutorial: [Get a LoadBalancer for your private Kubernetes cluster with inlets-operator](https://blog.alexellis.io/ingress-for-your-local-kubernetes-cluster/)

See also: [advanced usage of inlets including Docker, Kubernetes, multiple-services, and binding to private IPs](./docs/advanced.md)

### What are people saying about inlets?

Read community tutorials, the launch posts on Hacker News, and send a PR if you have written about inlets or inlets PRO:

* [Community and tutorials page](docs/community.md)

> You can share about inlets using `@inletsdev`, `#inletsdev`, and `https://inlets.dev`.

### Using inlets at work or in production?

See [ADOPTERS.md](./ADOPTERS.md) for what companies are doing with inlets today.

### Development

See [CONTRIBUTING.md](./CONTRIBUTING.md)

#### Other Kubernetes port-forwarding tooling

* [`kubectl port-forward`](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) - built into the Kubernetes CLI, forwards a single port to the local computer.
* [kubefwd](https://github.com/txn2/kubefwd) - Kubernetes utility to port-forward multiple services to your local computer.
* [kurun](https://github.com/banzaicloud/kurun) - Run main.go in Kubernetes with one command, also port-forward your app into Kubernetes.

inlets&reg; is a registered trademark of OpenFaaS Ltd. All rights reserved, registered company in the UK: 11076587