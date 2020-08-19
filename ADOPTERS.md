# Adopters

This is a list of adopters of inlets

## Adopters list (alphabetical)

- [Banzai Cloud](https://banzaicloud.com/) - inlets is used in [kurun](https://github.com/banzaicloud/kurun) for [proxying services](https://banzaicloud.com/blog/kurun-port-forward/).

- [Cloodot](https://beta.cloodot.com/) - We use inlets as a replacement for ngrok while developing with webhooks. This helps us develop integrations with facebook, twitter, whatsapp and google business.

- [VSHN](https://vshn.ch) - VSHN uses inlets for connecting managed OpenShift clusters
  which are behind a firewall to the customer portal. Only the control
  plane of inlets is exposed to the internet, the customer portal connects
  directly to the data plane which is exposed as a Kubernetes service in
  the same namespace as the customer portal. This makes sure no connections
  from outside the cluster on which the customer portal runs are able to
  connect the the data plane. Each managed OpenShift cluster has the inlets
  client installed, connecting to the central inlets server over HTTPS.
  The Managed OpenShift clusters are part of the [APPUiO](https://appuio.ch)
  product.
