# inlets

Expose your local endpoints to the Internet

[![Build Status](https://travis-ci.org/alexellis/inlets.svg?branch=master)](https://travis-ci.org/alexellis/inlets) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/alexellis/inlets)](https://goreportcard.com/report/github.com/alexellis/inlets) [![Documentation](https://godoc.org/github.com/alexellis/inlets?status.svg)](http://godoc.org/github.com/alexellis/inlets)

## Intro

inlets combines a reverse proxy and websocket tunnels to expose your internal and development endpoints to the public Internet via an exit-node. An exit-node may be a 5-10 USD VPS or any other computer with an IPv4 IP address.

Why do we need this project? Similar tools such as [ngrok](https://ngrok.com/) or [Argo Tunnel](https://developers.cloudflare.com/argo-tunnel/) from [Cloudflare](https://www.cloudflare.com/) are closed-source, have limits built-in, can work out expensive and have limited support for arm/arm64. Ngrok is also often banned by corporate firewall policies meaning it can be unusable. Other open-source tunnel tools are designed to only set up a static tunnel. inlets aims to dynamically bind and discover your local services to DNS entries with automated TLS certificates to a public IP address over its websocket tunnel.

When combined with SSL - inlets can be used with any corporate HTTP proxy which supports `CONNECT`.

![](docs/inlets.png)

*Conceptual diagram for inlets*

### Who is behind this project?

inlets is brought to you by [Alex Ellis](https://twitter.com/alexellisuk), the founder of [OpenFaaS](https://github.com/openfaas/faas/).

> The mission of [OpenFaaS](https://github.com/openfaas/faas/) is to Make Serverless Functions Simple for developers. With OpenFaaS you can package any code, binary or microservice into a Serverless Function and deploy to Kubernetes or Docker Swarm without repetitive boiler-plate coding or complex YAML files. OpenFaaS has over 15.5k GitHub stars, 180 contributors and a growing end-user community.

Make a donation to support inlets via [PayPal](https://www.paypal.me/openfaas?locale.x=en_GB).

### Goals

#### Initial goals:

* automatically create endpoints on exit-node based upon client definitions
  * multiplex sites on same port through use of DNS / host entries 
* link encryption using SSL over websockets (`wss://`)
* automatic reconnect
* authentication using service account or basic auth
* automatic TLS provisioning for endpoints using [cert-magic](https://github.com/mholt/certmagic)
  * configure staging or production LetsEncrypt issuer using HTTP01 challenge

#### Stretch goals:

* discover and configure endpoints for Ingress definitions from Kubernetes
* configuration to run "exit-node" as serverless container with Azure ACI / AWS Fargate
* automatic configuration of DNS / A records
* configure staging or production LetsEncrypt issuer using DNS01 challenge

#### Non-goals:

* tunnelling plain (non-HTTP) traffic over TCP

### Status

Unlike HTTP 1.1 which follows a synchronous request/response model websockets use an asynchronous pub/sub model for sending and receiving messages. This presents a challenge for tunneling a synchronous protocol over an asynchronous bus. This is a working prototype that can be used for testing, development and to generate discussion, but is not production-ready.

* ~~There is currently no authentication on the server component~~ The tunnel link is secured via `-token` flag and a shared secret
* The default configuration uses websockets without SSL `ws://`, but to enable encryption you could enable SSL `wss://`
* ~~There is no timeout for when the tunnel is disconnected~~ timeout can be configured via args on the server
* ~~The upstream URL has to be configured on both server and client until a discovery or service advertisement mechanism is added~~ advertise on the client

### Get `inlets`

Run this command:

```bash
# Install to local directory
curl -sLS https://get.inlets.dev | sh

# Install to /usr/local/bin/
curl -sLS https://get.inlets.dev | sudo sh
```

Binaries for Linux, Darwin (MacOS) and armhf are made available via the [releases page](https://github.com/alexellis/inlets/releases). You will also find SHA checksums available if you want to verify your download.

## Test it out

* On the server or exit-node

Start the tunnel server on a machine with a publicly-accessible IPv4 IP address such as a VPS.

```bash
./inlets server --port=80
```

> Note: You can pass the `-token` argument followed by a token value to both the server and client to prevent unauthorized connections to the tunnel.

Example with token:

```bash
token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1); ./inlets server --port=8090 --token="$token"
```

Note down your public IPv4 IP address i.e. 192.168.0.101

* On your machine behind the firewall start an example service that you want to expose to the Internet

You can use my hash-browns service for instance which generates hashes.

Install hash-browns or run your own HTTP server

```sh
go get -u github.com/alexellis/hash-browns
cd $GOPATH/src/github.com/alexellis/hash-browns

port=3000 go run server.go 
```

* On your machine behind the firewall

Start the tunnel client

```sh
./inlets client \
 --remote=192.168.0.101:80 \
 --upstream=http://127.0.0.1:3000
```

Replace the `--remote` with the address where your other machine is listening.

We now have an example service running (hash-browns), a tunnel server and a tunnel client.

So send a request to the public IP address or hostname:

```sh 
./inlets client --remote=192.168.0.101:80 --upstream  "gateway.mydomain.tk=http://127.0.0.1:3000"
```

```sh
curl -d "hash this" http://192.168.0.101/hash -H "Host: gateway.mydomain.tk"
# or
curl -d "hash this" http://192.168.0.101/hash
# or
curl -d "hash this" http://gateway.mydomain.tk/hash
```

You will see the traffic pass between the exit node / server and your development machine. You'll see the hash message appear in the logs as below:

```sh
~/go/src/github.com/alexellis/hash-browns$ port=3000 go run server.go 
2018/12/23 20:15:00 Listening on port: 3000
"hash this"
```

Now check the metrics endpoint which is built-into the hash-browns example service:

```sh
curl http://192.168.0.101/metrics | grep hash
```

## Development

For development you will need Golang 1.10 or 1.11 on both the exit-node or server and the client.

You can get the code like this:

```bash
go get -u github.com/alexellis/inlets
cd $GOPATH/src/github.com/alexellis/inlets
```

Contributions are welcome. All commits must be signed-off with `git commit -s` to accept the [Developer Certificate of Origin](https://developercertificate.org).

## Take things further

You can expose an OpenFaaS or OpenFaaS Cloud deployment with `inlets` - just change `--upstream=http://127.0.0.1:3000` to `--upstream=http://127.0.0.1:8080` or `--upstream=http://127.0.0.1:31112`. You can even point at an IP address inside or outside your network for instance: `--upstream=http://192.168.0.101:8080`.

You can build a basic supervisor script for `inlets` in case of a crash, it will re-connect within 5 seconds:

In this example the Host/Client is acting as a relay for OpenFaaS running on port 8080 on the IP 192.168.0.28 within the internal network.

Host/Client:

```sh
while [ true ] ; do sleep 5 && ./inlets client --upstream=http://192.168.0.28:8080 --remote=exit.my.club  ; done
```

Exit-node:

```sh
while [ true ] ; do sleep 5 && ./inlets server --upstream=http://192.168.0.28:8080 ; done
```

### Run as a deployment on Kubernetes

You can even run `inlets` within your Kubernetes in Docker (kind) cluster to get ingress (incoming network) for your services such as the OpenFaaS gateway:

```yaml
apiVersion: apps/v1beta1 # for versions before 1.6.0 use extensions/v1beta1
kind: Deployment
metadata:
  name: inlets
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: inlets
    spec:
      containers:
      - name: inlets
        image: alexellis2/inlets:0.5.4
        imagePullPolicy: Always
        command: ["./inlets"]
        args:
        - "client"
        - "--upstream=http://gateway.openfaas:8080,http://endpoint.openfaas:9090"
        - "--remote=your-public-ip"
```

Replace the line: `- "--remote=your-public-ip"` with the public IP belonging to your VPS.

### Run on a VPS

Provisioning on a VPS will see inlets running as a systemd service.  All the usual `service` commands should be used with `inlets` as the service name.

Inlets uses a token to prevent unauthorized access to the server component.  A known token can be configured by amending [userdata.sh](./hack/userdata.sh) prior to provisioning

```sh
# Enables randomly generated authentication token by default.
# Change the value here if you desire a specific token value.
export INLETSTOKEN=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
```

If the token value is randomly generated then you will need to access the VPS in order to obtain the token value.

```sh
cat /etc/default/inlets 
```  

#### Where can I get a cheap / free domain-name?

You can get a free domain-name with a .tk / .ml or .ga TLD from https://www.freenom.com - make sure the domain has at least 4 letters to get it for free. You can also get various other domains starting as cheap as 1-2USD from https://www.namecheap.com

[Namecheap](https://www.namecheap.com) provides wildcard TLS out of the box, but [freenom](https://www.freenom.com) only provides root/naked domain and a list of sub-domains. Domains from both providers can be moved to alternative nameservers for use with AWS Route 53 or Google Cloud DNS - this then enables wildcard DNS and the ability to get a wildcard TLS certificate from LetsEncrypt.

#### Digital Ocean

If you are a Digital Ocean user and use `doctl` then you can provision a host with [./hack/provision-digitalocean.sh](./hack/provision-digitalocean.sh).  Please ensure you have configured `droplet.create.ssh-keys` within your `~/.config/doctl/config.yaml`.

#### Scaleway

[Scaleway](https://www.scaleway.com/) offer probably the cheapest option at 1.99 EUR / month using the "1-XS" from the "Start" tier. 

If you have the Scaleway CLI installed you can provision a host with [./hack/provision-scaleway.sh](./hack/provision-scaleway.sh).
