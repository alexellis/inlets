## Quickstart tutorial

This tutorial shows you how the inlets client and server components work together, using your laptop to run both parts of the tunnel.

> See also: [how to provision an "exit-server" with a public IPv4 address using a VPS](./vps.md).

* On the *exit-server* (or server)

Start the tunnel server on a machine with a publicly-accessible IPv4 IP address such as a VPS.

Example with a token for client authentication:

```bash
export token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
inlets server --port=8090 --token=$token
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
export GO111MODULE=off
export GOPATH=$HOME/go/

go get -u github.com/alexellis/hash-browns
port=3000 $GOPATH/bin/hash-browns
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
 --token=$TOKEN
```

* Replace the `--remote` with the address where your exit-server is running `inlets server`.
* Replace the `--token` with the value from your server

We now have three processes:
* example service running (hash-browns) or Python's webserver
* an exit-server running the tunnel server (`inlets server`)
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

You will see the traffic pass between the exit server / server and your development machine. You'll see the hash message appear in the logs as below:

```sh
~/go/src/github.com/alexellis/hash-browns$ port=3000 go run main.go
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
 --upstream="store1.example.com=http://127.0.0.1:8001,store2.example.com=http://127.0.0.1:8002" \
 --token=$TOKEN
```

You can now create two DNS entries or `/etc/hosts` file entries for `store1.example.com` and `store2.example.com`, then connect through your browser.
