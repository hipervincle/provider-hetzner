<!--
SPDX-FileCopyrightText: 2023 The Crossplane Authors <https://crossplane.io>

SPDX-License-Identifier: CC-BY-4.0
-->

# Upjet-based Crossplane provider for Hetzner

<div style="text-align: center;">

![CI](https://github.com/miaits/provider-hetzner/workflows/CI/badge.svg)
[![GitHub release](https://img.shields.io/github/release/miaits/provider-hetzner/all.svg)](https://github.com/miaits/provider-hetzner/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/miaits/provider-hetzner)](https://goreportcard.com/report/github.com/miaits/provider-hetzner)
[![Contributors](https://img.shields.io/github/contributors/miaits/provider-hetzner)](https://github.com/miaits/provider-hetzner/graphs/contributors)

</div>

Provider Hetzner is a [Crossplane](https://crossplane.io/) provider that is
built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for
[Hetzner Cloud](https://www.hetzner.com/cloud).

## Getting Started

The provider needs a Kubernetes secret with an API token
for Hetzner Cloud. Create it in `crossplane-system`:

```bash
kubectl create secret generic hetzner   \
-n crossplane-system                    \
--from-literal=credentials='{"token":"<TOKEN>"}'   \
--dry-run=client -o yaml | kubectl apply -f -
```

To install the provider into a local Kubernetes cluster with 
Crossplane already installed, apply:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-hetzner
spec:
  package: ghcr.io/hipervincle/provider-hetzner:v0.2.1
```

After installing the provider, create a provider configuration before
applying managed resources:


```yaml
apiVersion: hetzner.m.crossplane.io/v1beta1
kind: ClusterProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: hetzner
      namespace: crossplane-system
      key: credentials
```

To create a Hetzner server for test, apply:

```yaml
apiVersion: server.hetzner.m.crossplane.io/v1alpha1
kind: Server
metadata:
  annotations:
    meta.upbound.io/example-id: server/v1alpha1/server
  labels:
    testing.upbound.io/example-name: node1
  name: node1
spec:
  forProvider:
    image: debian-12
    name: node1
    serverType: cx23
```
To delete the test server, run:
```bash
kubectl delete servers.server.hetzner.m node1
```

Reference-first namespaced examples are available at:

- `examples/namespaced/server/v1alpha1/simple-server.yaml`
- `examples/namespaced/server/v1alpha1/networkattachment.yaml`
- `examples/namespaced/server/v1alpha1/firewallattachment.yaml`
- `examples/namespaced/loadbalancer/v1alpha1/networkattachment.yaml`
- `examples/namespaced/loadbalancer/v1alpha1/loadbalancerservice.yaml`
- `examples/web-stack/web-stack-foundation.yaml`
- `examples/web-stack/web-stack.yaml`
- `examples/namespaced/dns/zone-and-records.yaml`

DNS zone and RRSet implementation details, generation flow, runtime reconciliation,
usage, and import guidance are documented in
[`docs/adding-dns-zone-management.md`](docs/adding-dns-zone-management.md).

The private-backend web stack is split into two phases:

1. `kubectl apply -f examples/web-stack/web-stack-foundation.yaml`
2. Wait for the subnet, gateway, routes, and load balancer to become ready.
3. `kubectl apply -f examples/web-stack/web-stack.yaml`

The main renamed kinds in `v1alpha1` are `LoadBalancer`, `NetworkAttachment`,
`PlacementGroup`, `SSHKey`, `FloatingIP`, and `FloatingIPAssignment`.

## Contributing

For the general contribution guide,
see [Upjet Contribution Guide](https://github.com/crossplane/upjet/blob/main/CONTRIBUTING.md)

If you'd like to learn how to use Upjet, see [Usage Guide](https://github.com/crossplane/upjet/tree/main/docs).

To build this provider locally and run it in a local Kubernetes cluster, run
`make run`.

The local `make run` target binds metrics and webhook listeners to ephemeral
ports by default to avoid conflicts with other local processes. You can
override them when needed, for example:
`make run LOCAL_METRICS_BIND_ADDRESS=127.0.0.1:8081 LOCAL_WEBHOOK_PORT=9444`

### Add a New Resource

Follow the Upjet guide
for [adding new resources](https://github.com/crossplane/upjet/blob/main/docs/adding-new-resource.md).

## Getting help

For filing bugs, suggesting improvements, or requesting new resources or features, please
open an [issue](https://github.com/miaits/provider-hetzner/issues/new/choose).


## License

The provider is released under the [the Apache 2.0 license](LICENSE) with [notice](NOTICE).
