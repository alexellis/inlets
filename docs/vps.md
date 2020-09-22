# Run Inlets on a VPS

## The easy way - `inletsctl`

`inletsctl` can provision an exit-server along with your systemd unit file all within a few seconds.

```bash
inletsctl create --provider digitalocean \
 --region lon1 \
 --access-token-file $HOME/do-token
```

Then connect with the parameters printed out for you.

See more in the [inletsctl reference documentation](https://docs.inlets.dev/#/tools/inletsctl?id=inletsctl-reference-documentation)

## The long way

Provisioning on a VPS will see inlets running as a systemd service.  All the usual `service` commands should be used with `inlets` as the service name.

Inlets uses a token to prevent unauthorized access to the server component.  A known token can be configured by amending [userdata.sh](../hack/userdata.sh) prior to provisioning

```sh
# Enables randomly generated authentication token by default.
# Change the value here if you desire a specific token value.
export INLETSTOKEN=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
```

If the token value is randomly generated then you will need to access the VPS in order to obtain the token value.

```sh
cat /etc/default/inlets
```

## How do I enable TLS / HTTPS?

* Create a DNS A record for your exit-server IP and the DNS entry `exit.domain.com` (replace as necessary).

* Download Caddy from the [Releases page](https://github.com/mholt/caddy/releases).

* Enter this text into a Caddyfile replacing `exit.domain.com` with your subdomain.

### Caddy v1
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

### Caddy v2
```Caddyfile
exit.domain.com

reverse_proxy / 127.0.0.1:8000 {
  header_up Host {http.request.host}
  header_up X-Real-IP {http.request.remote}
  header_up X-Forwarded-For {http.request.remote}
  header_up X-Forwarded-Port {http.request.port}
  header_up X-Forwarded-Proto {http.request.scheme}
}

reverse_proxy /tunnel 127.0.0.1:8000 {
  header_up Host {http.request.host}
  header_up X-Real-IP {http.request.remote}
  header_up X-Forwarded-For {http.request.remote}
  header_up X-Forwarded-Port {http.request.port}
  header_up X-Forwarded-Proto {http.request.scheme}
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

You now have a secure TLS link between your client(s) and your server on the exit server and for your site to serve traffic over.

## Where can I get a cheap / free domain-name?

You can get a free domain-name with a .tk / .ml or .ga TLD from https://www.freenom.com - make sure the domain has at least 4 letters to get it for free. You can also get various other domains starting as cheap as 1-2USD from https://www.namecheap.com

[Namecheap](https://www.namecheap.com) provides wildcard TLS out of the box, but [freenom](https://www.freenom.com) only provides root/naked domain and a list of sub-domains. Domains from both providers can be moved to alternative nameservers for use with AWS Route 53 or Google Cloud DNS - this then enables wildcard DNS and the ability to get a wildcard TLS certificate from LetsEncrypt.

My recommendation: pay to use [Namecheap](https://www.namecheap.com).

## Where can I host an `inlets` exit-server?

You can use inlets to provide incoming connections to any network, including containers, VM and AWS Firecracker.

Examples:

* Green to green - from one internal LAN to another
* Green to red - from an internal network to the Internet (i.e. Raspberry Pi cluster)
* Red to green - to make a service on a public network accessible as if it were a local service.

The following VPS providers have credit, or provisioning scripts to get an exit-server in a few moments.

Installation scripts have been provided which use `systemd` as a process supervisor. This means that if inlets crashes, it will be restarted automatically and logs are available.

* After installation, find your token with `sudo cat /etc/default/inlets`
* Check logs with `sudo systemctl status inlets`
* Restart with `sudo systemctl restart inlets`
* Check config with `sudo systemctl cat inlets`

### DigitalOcean

If you're a [DigitalOcean](https://www.digitalocean.com) user and use `doctl` then you can provision a host with [../hack/provision-digitalocean.sh](../hack/provision-digitalocean.sh).  Please ensure you have configured `droplet.create.ssh-keys` within your `~/.config/doctl/config.yaml`.

DigitalOcean will then email you the IP and root password for your new host. You can use it to log in and get your auth token, so that you can connect your client after that.

Datacenters for exit-servers are available world-wide

### Civo

[Civo](https://www.civo.com/) is a UK developer cloud and [offers 50 USD free credit](http://bit.ly/2Lx9d2o).

Installation is currently manual and the datacenter is located in London.

* Create a VM of any size and then download and run inlets as a server
* Copy over `../hack/userdata.sh` and run it on the server as `root`

### Scaleway

[Scaleway](https://www.scaleway.com/) offer probably the cheapest option at 1.99 EUR / month using the "1-XS" from the "Start" tier.

If you have the Scaleway CLI installed you can provision a host with [../hack/provision-scaleway.sh](../hack/provision-scaleway.sh).

Datacenters include: Paris and Amsterdam.

### Running over an SSH tunnel

You can tunnel over SSH if you are not using a reverse proxy that enables SSL. This encrypts the traffic over the tunnel.

On your client, create a tunnel to the exit-server:

```
ssh -L 8000:127.0.0.1:80 exit-server-ip
```

Now for the `--remote` address use `--remote ws://127.0.0.1:8000`

