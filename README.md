# inlets

Expose your local endpoints to the Internet

[![Build Status](https://travis-ci.org/alexellis/inlets.svg?branch=master)](https://travis-ci.org/alexellis/inlets) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/alexellis/inlets)](https://goreportcard.com/report/github.com/alexellis/inlets) [![Documentation](https://godoc.org/github.com/alexellis/inlets?status.svg)](http://godoc.org/github.com/alexellis/inlets)

## Intro

inlets combines a reverse proxy and websocket tunnels to expose your internal and development endpoints to the public Internet via an exit-node. An exit-node may be a 5-10 USD VPS or any other computer with an IPv4 IP address.

Why do we need this project? Similar tools such as [ngrok](https://ngrok.com/) or [Argo Tunnel](https://developers.cloudflare.com/argo-tunnel/) from [Cloudflare](https://www.cloudflare.com/) are closed-source, have limits built-in, can work out expensive, and have limited support for arm/arm64. Ngrok is also often banned by corporate firewall policies meaning it can be unusable. Other open-source tunnel tools are designed to only set up a single static tunnel. inlets aims to dynamically bind and discover your local services to DNS entries with automated TLS certificates to a public IP address over a websocket tunnel.

When combined with SSL - inlets can be used with any corporate HTTP proxy which supports `CONNECT`.

![](docs/inlets.png)

*Conceptual diagram for inlets*

## License & terms

**Important**

Developers wishing to use inlets within a corporate network are advised to seek approval from their administrators or management before using the tool. By downloading, using, or distributing inlets, you agree to the [LICENSE](./LICENSE) terms & conditions. No warranty or liability is provided.

### Who is behind this project?

inlets is brought to you by [Alex Ellis](https://twitter.com/alexellisuk). Alex is a [CNCF Ambassador](https://www.cncf.io/people/ambassadors/) and the founder of [OpenFaaS](https://github.com/openfaas/faas/).

> The mission of [OpenFaaS](https://github.com/openfaas/faas/) is to Make Serverless Functions Simple for developers. With OpenFaaS you can package any code, binary or microservice into a Serverless Function and deploy to Kubernetes or Docker without repetitive boiler-plate coding or complex YAML files. OpenFaaS has over 17.5k GitHub stars, 200 contributors, and a growing end-user community.

Become an Insider to receive regular *Insider Updates* on inlets, and all his other OSS work, blogs and videos via [GitHub Sponsors](https://github.com/users/alexellis/sponsorship)

### Goals

#### Initial goals:

* automatically create endpoints on exit-node based upon client definitions
  * multiplex sites on same port and websocket through the use of DNS / host entries
* link encryption using SSL over websockets (`wss://`)
* automatic reconnect
* authentication using service account or basic auth
* automatic TLS provisioning for endpoints using [cert-magic](https://github.com/mholt/certmagic)
  * configure staging or production LetsEncrypt issuer using HTTP01 challenge
* native multi-arch with ARMHF/ARM64 support
* integration with Docker and [Kubernetes](https://kubernetes.io)

#### Stretch goals:

* automatic configuration of DNS / A records
* tunnelling websocket traffic in addition to HTTP(s)
* discover and configure endpoints for `Ingress` definitions from [Kubernetes](https://kubernetes.io)
* configuration to run "exit-node" as serverless container with Azure ACI / AWS Fargate
* configure staging or production LetsEncrypt issuer using DNS01 challenge
* [get a logo for the project](https://github.com/alexellis/inlets/issues/46)

#### Non-goals:

* tunnelling plain TCP traffic over the websocket

> Note: this is a non-goal at present, but I am open to contributions

### Status

Unlike HTTP 1.1 which follows a synchronous request/response model websockets use an asynchronous pub/sub model for sending and receiving messages. This presents a challenge for tunneling a *synchronous protocol* over an *asynchronous bus*.

inlets 2.0 introduces performance enhancements and leverages parts of the Kubernetes and Rancher API. It uses the same tunnelling packages that enable node-to-node communication in [Rancher's k3s project](https://k3s.io). It is suitable for development and may be useful in production. Before deploying `inlets` into production, it is advised that you do adequate testing.

Feel free to open issues if you have comments, suggestions or contributions.

* The tunnel link is secured via `--token` flag using a shared secret
* The default configuration uses websockets without SSL `ws://`, but to enable encryption you could enable SSL `wss://`
* A timeout for requests can be configured via args on the server
* ~~The upstream URL has to be configured on both server and client until a discovery or service advertisement mechanism is added~~ The client can advertise upstream URLs, which it can serve
* The tunnel transport is wrapped by default which strips CORS headers from responses, but you can disable it with the `--disable-transport-wrapping` flag on the server

### Video demo

Using inlets I was able to set up a public endpoint (with a custom domain name) for my JavaScript & Webpack [Create React App](https://github.com/facebook/create-react-app).

[![https://img.youtube.com/vi/jrAqqe8N3q4/hqdefault.jpg](https://img.youtube.com/vi/jrAqqe8N3q4/maxresdefault.jpg)](https://youtu.be/jrAqqe8N3q4)

### What are people saying about inlets?

> You can share about inlets using `#inletsdev`, `#inlets`, and `https://inlets.dev`.

inlets has trended on the front page of Hacker News twice.

* [inlets 1.0](https://news.ycombinator.com/item?id=19189455) - 146 points, 48 comments
* [inlets 2.0](https://news.ycombinator.com/item?id=20410552) - 218 points, 66 comments

Twitter:

* ["I just transferred a 70Gb disk image from a NATed NAS to a remote NATed server with @alexellisuk inlets tunnels and a one-liner python web server" by Roman Dodin](https://twitter.com/ntdvps/status/1143071544203186176)
* ["Really amazed by inlets by @alexellisuk - "Up and running in 15min - I will be able to watch my #RaspberryPi servers running at home while staying on the beach 🏄‍♂️🌴🍸👏👏👏" by Florian Dambrine](https://twitter.com/DambrineF/status/1158364581624012802?s=20)
* [Testing an OAuth proxy by Vivek Singh](https://twitter.com/viveksyngh/status/1142054203478564864)
* [inlets used at KubeCon to power a live IoT demo at a booth](https://twitter.com/tobruzh/status/1130421702914129921)
* [PR to support Risc-V by Carlos Eduardo](https://twitter.com/carlosedp/status/1140740494617645061)
* [Recommended by Michael Hausenblas for use with local Kubernetes](https://twitter.com/mhausenblas/status/1143020953380753409)
* [5 top facts about inlets by Alex Ellis](https://twitter.com/alexellisuk/status/1140552115204608001)
* ["Cool! I hadn't heard of inlets until now, but I love the idea of exposing internal services this way. I've been using TOR to do this!" by Stephen Doskett, Tech Field Day](https://twitter.com/SFoskett/status/1108989190912524288)
* ["Learn how to set up HTTPS for your local endpoints with inlets, Caddy, and DigitalOcean thanks to @alexellisuk!" by @DigitalOcean](https://twitter.com/digitalocean/status/1113440166310502400)
* ["See how Inlets helped me to expose my local endpoints for my homelab that sits behind a Carrier-Grade NAT"](https://twitter.com/ruanbekker/status/1161399537417801728)
* ["Securely access external applications as Kubernetes Services, from your laptop or from any other host, using inlets"](https://twitter.com/BanzaiCloud/status/1164168218954670080)

> Note: add a PR to send your story or use-case, I'd love to hear from you.


## Get started

You can install the CLI with a `curl` utility script, `brew` or by downloading the binary from the releases page. Once installed you'll get the `inlets` command.

### Install the CLI

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

Binaries are made available on the [releases page](https://github.com/alexellis/inlets/releases) for Linux (x86_64, armhf & arm64) and for Darwin (MacOS). You will also find SHA checksums available if you want to verify your download.

### Test it out

You can run inlets between any two computers with connectivity, these could be containers, VMs, bare metal or even "loop-back" on your own laptop.

See [how to provision an "exit-node" with a public IPv4 address using a VPS](#run-on-a-vps).

* On the *exit-node* (or server)

Start the tunnel server on a machine with a publicly-accessible IPv4 IP address such as a VPS.

Example with a token for client authentication:

```bash
export token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
inlets server --port=8090 --token="$token"
```

> Note: You can pass the `--token` argument followed by a token value to both the server and client to prevent unauthorized connections to the tunnel.


```bash
inlets server --port=8090
```

You can also run unprotected, but this is not recommended.

Note down your public IPv4 IP address.

* Head over to your machine where you are running a sample service, or something you want to expose.

You can use my hash-browns service for instance which generates hashes.

Install hash-browns or run your own HTTP server

```sh
go get -u github.com/alexellis/hash-browns
cd $GOPATH/src/github.com/alexellis/hash-browns

port=3000 go run server.go
```

If you don't have Go installed, then you could run [Python's built-in HTTP server](https://docs.python.org/2/library/simplehttpserver.html):

```sh
mkdir -p /tmp/inlets-test/
cd /tmp/inlets-test/
touch hello-world
python -m SimpleHTTPServer 3000
```

* On the same machine, start the inlets client

Start the tunnel client:

```sh
export REMOTE="127.0.0.1:8090"    # for testing inlets on your laptop, replace with the public IPv4
export TOKEN="CLIENT-TOKEN-HERE"  # the client token is found on your VPS or on start-up of "inlets server"
inlets client \
 --remote=$REMOTE \
 --upstream=http://127.0.0.1:3000 \
 --token $TOKEN
```

* Replace the `--remote` with the address where your exit-node is running `inlets server`.
* Replace the `--token` with the value from your server

We now have three processes:
* example service running (hash-browns) or Python's webserver
* an exit-node running the tunnel server (`inlets server`)
* a client running the tunnel client (`inlets client`)

So send a request to the inlets server - use its domain name or IP address:

Assuming `gateway.mydomain.tk` points to `127.0.0.1` in `/etc/hosts` or your DNS server.

```sh
curl -d "hash this" http://127.0.0.1:8090/hash -H "Host: gateway.mydomain.tk"
# or
curl -d "hash this" http://127.0.0.1:8090/hash
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
curl $REMOTE/metrics | grep hash
```

You can also use multiple domain names and tie them back to different internal services.

Here we start the Python server on two different ports, serving content from two different locations and then map it to two different Host headers, or domain names:

```sh
mkdir -p /tmp/store1
cd /tmp/store1/
touch hello-store-1
python -m SimpleHTTPServer 8001 &


mkdir -p /tmp/store2
cd /tmp/store2/
touch hello-store-2
python -m SimpleHTTPServer 8002 &
```

```sh
export REMOTE="127.0.0.1:8090"    # for testing inlets on your laptop, replace with the public IPv4
export TOKEN="CLIENT-TOKEN-HERE"  # the client token is found on your VPS or on start-up of "inlets server"
inlets client \
 --remote=$REMOTE \
 --token $TOKEN \
 --upstream="store1.example.com=http://127.0.0.1:8001,store2.example.com=http://127.0.0.1:8002"
```

You can now create two DNS entries or `/etc/hosts` file entries for `store1.example.com` and `store2.example.com`, then connet through your browser.

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
while [ true ] ; do sleep 5 && inlets client --upstream=http://192.168.0.28:8080 --remote=exit.my.club  ; done
```

Exit-node:

```sh
while [ true ] ; do sleep 5 && inlets server --upstream=http://192.168.0.28:8080 ; done
```

## Bind a different port for the control-plane

You can bind two separate TCP ports for the user-facing port and the tunnel.

* `--port` - the port for users to connect to and for serving data, i.e. the *Data Plane*
* `--control-plane` - the port for the websocket to connect to i.e. the *Control Plane*

### Docker & Kubernetes application development

Docker images are published for `x86_64` and `armhf`

* `alexellis2/inlets:2.3.2`
* `alexellis2/inlets:2.3.2-armhf`

> Note: For Raspberry Pi, you need to use the image ending in `-armhf`.

#### Run as a deployment on Kubernetes

You can run the client inside Kubernetes to expose your local services to the Internet, or another network.

Here's an example showing how to get ingress into your cluster for your OpenFaaS gateway and for Prometheus:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inlets
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: inlets
  template:
    metadata:
      labels:
        app.kubernetes.io/name: inlets
    spec:
      containers:
      - name: inlets
        image: alexellis2/inlets:2.3.2
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "client"
        - "--upstream=http://gateway.openfaas:8080,http://prometheus.openfaas:9090"
        - "--remote=your-public-ip"
```

Replace the line: `- "--remote=your-public-ip"` with the public IP belonging to your VPS.

Alternatively, see the unofficial helm chart from the community: [inlets-helm](https://github.com/paurosello/inlets_helm).

> Note: For Raspberry Pi, you need to use the image ending in `-armhf`.

#### Use authentication from a Kubernetes secret

In production, you should always use a secret to protect your exit-node. You will need a way of passing that to your server and inlets allows you to read a Kubernetes secret.

* Create a random secret

```
$ kubectl create secret generic inlets-token --from-literal token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
secret/inlets-token created
```

* Or create a secret with the value from your remote server

```
$ export TOKEN=""
$ kubectl create secret generic inlets-token --from-literal token=${TOKEN}
secret/inlets-token created
```

* Bind the secret named `inlets-token` to the Deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inlets
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: inlets
  template:
    metadata:
      labels:
        app.kubernetes.io/name: inlets
    spec:
      containers:
      - name: inlets
        image: alexellis2/inlets:2.3.2
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "client"
        - "--remote=ws://REMOTE-IP"
        - "--upstream=http://gateway.openfaas:8080"
        - "--token-from=/var/inlets/token"
        volumeMounts:
          - name: inlets-token-volume
            mountPath: /var/inlets/
      volumes:
        - name: inlets-token-volume
          secret:
            secretName: inlets-token
```

Optional tear-down:

```
$ kubectl delete deploy/inlets
$ kubectl delete secret/inlets-token
```

#### Use your Kubernetes cluster as an exit-node

You can use a Kubernetes cluster which has public IP addresses, an IngressController, or a LoadBalancer to run one or more exit-nodes.

* Create a random secret

```
$ kubectl create secret generic inlets-token --from-literal token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
secret/inlets-token created
```

* Or create a secret with the value from your remote server

```
$ export TOKEN=""
$ kubectl create secret generic inlets-token --from-literal token=${TOKEN}
secret/inlets-token created
```

* Create a `Service`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: inlets
  labels:
    app: inlets
spec:
  type: ClusterIP
  ports:
    - port: 8000
      protocol: TCP
      targetPort: 8000
  selector:
    app: inlets
```

* Create a `Deployment`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inlets
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: inlets
  template:
    metadata:
      labels:
        app.kubernetes.io/name: inlets
    spec:
      containers:
      - name: inlets
        image: alexellis2/inlets:2.3.2
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "server"
        - "--token-from=/var/inlets/token"
        volumeMounts:
          - name: inlets-token-volume
            mountPath: /var/inlets/
      volumes:
        - name: inlets-token-volume
          secret:
            secretName: inlets-token
```

You can now create an `Ingress` record, or `LoadBalancer` to connect to your server.

#### Try inlets with KinD (Kubernetes in Docker)

Try this guide to expose services running in a [KinD cluster](https://github.com/kubernetes-sigs/kind):

[Micro-tutorial inlets with KinD](https://gist.github.com/alexellis/c29dd9f1e1326618f723970185195963)

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

#### How do I enable TLS / HTTPS?

* Create a DNS A record for your exit-node IP and the DNS entry `exit.domain.com` (replace as necessary).

* Download Caddy from the [Releases page](https://github.com/mholt/caddy/releases).

* Enter this text into a Caddyfile replacing `exit.domain.com` with your subdomain.

```Caddyfile
exit.domain.com

proxy / 127.0.0.1:8000 {
  transparent
}

proxy /tunnel 127.0.0.1:8000 {
  transparent
  websocket
}
```

* Run `inlets server --port 8000`

* Run `caddy`

Caddy will now ask you for your email address and after that will obtain a TLS certificate for you.

* On the client run the following, adding any other parameters you need for `--upstream`

```
inlets client --remote wss://exit.domain.com
```

> Note: wss indicates to use port 443 for TLS.

You now have a secure TLS link between your client(s) and your server on the exit node and for your site to serve traffic over.

#### Where can I get a cheap / free domain-name?

You can get a free domain-name with a .tk / .ml or .ga TLD from https://www.freenom.com - make sure the domain has at least 4 letters to get it for free. You can also get various other domains starting as cheap as 1-2USD from https://www.namecheap.com

[Namecheap](https://www.namecheap.com) provides wildcard TLS out of the box, but [freenom](https://www.freenom.com) only provides root/naked domain and a list of sub-domains. Domains from both providers can be moved to alternative nameservers for use with AWS Route 53 or Google Cloud DNS - this then enables wildcard DNS and the ability to get a wildcard TLS certificate from LetsEncrypt.

My recommendation: pay to use [Namecheap](https://www.namecheap.com).

#### Where can I host an `inlets` exit-node?

You can use inlets to provide incoming connections to any network, including containers, VM and AWS Firecracker.

Examples:

* Green to green - from one internal LAN to another
* Green to red - from an internal network to the Internet (i.e. Raspberry Pi cluster)
* Red to green - to make a service on a public network accessible as if it were a local service.

The following VPS providers have credit, or provisioning scripts to get an exit-node in a few moments.

Installation scripts have been provided which use `systemd` as a process supervisor. This means that if inlets crashes, it will be restarted automatically and logs are available.

* After installation, find your token with `sudo cat /etc/default/inlets`
* Check logs with `sudo systemctl status inlets`
* Restart with `sudo systemctl restart inlets`
* Check config with `sudo systemctl cat inlets`

##### DigitalOcean

If you're a [DigitalOcean](https://www.digitalocean.com) user and use `doctl` then you can provision a host with [./hack/provision-digitalocean.sh](./hack/provision-digitalocean.sh).  Please ensure you have configured `droplet.create.ssh-keys` within your `~/.config/doctl/config.yaml`.

DigitalOcean will then email you the IP and root password for your new host. You can use it to log in and get your auth token, so that you can connect your client after that.

Datacenters for exit-nodes are available world-wide

##### Civo

[Civo](https://www.civo.com/) is a UK developer cloud and [offers 50 USD free credit](http://bit.ly/2Lx9d2o).

Installation is currently manual and the datacenter is located in London.

* Create a VM of any size and then download and run inlets as a server
* Copy over `./hack/userdata.sh` and run it on the server as `root`

##### Scaleway

[Scaleway](https://www.scaleway.com/) offer probably the cheapest option at 1.99 EUR / month using the "1-XS" from the "Start" tier.

If you have the Scaleway CLI installed you can provision a host with [./hack/provision-scaleway.sh](./hack/provision-scaleway.sh).

Datacenters include: Paris and Amsterdam.

#### Running over an SSH tunnel

You can tunnel over SSH if you are not using a reverse proxy that enables SSL. This encrypts the traffic over the tunnel.

On your client, create a tunnel to the exit-node:

```
ssh -L 8000:127.0.0.1:80 exit-node-ip
```

Now for the `--remote` address use `--remote ws://127.0.0.1:8000`
