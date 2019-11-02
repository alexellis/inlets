# Inlets & Kubernetes recipes

## Run as a deployment on Kubernetes

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
        image: inlets/inlets:2.6.1
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "client"
        - "--upstream=http://gateway.openfaas:8080,http://prometheus.openfaas:9090"
        - "--remote=your-public-ip"
```

Replace the line: `- "--remote=your-public-ip"` with the public IP belonging to your VPS.

Alternatively, see the unofficial helm chart from the community: [inlets-helm](https://github.com/paurosello/inlets_helm).

## Use authentication from a Kubernetes secret

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
        image: inlets/inlets:2.6.1
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

## Use your Kubernetes cluster as an exit-node

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
    app.kubernetes.io/name: inlets
spec:
  type: ClusterIP
  ports:
    - port: 8000
      protocol: TCP
      targetPort: 8000
  selector:
    app.kubernetes.io/name: inlets
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
        image: inlets/inlets:2.6.1
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
Note that clients connecting to this server will have to specify port 8000 for their remote, as the default is 80.

## Try inlets with KinD (Kubernetes in Docker)

Try this guide to expose services running in a [KinD cluster](https://github.com/kubernetes-sigs/kind):

[Micro-tutorial inlets with KinD](https://gist.github.com/alexellis/c29dd9f1e1326618f723970185195963)
