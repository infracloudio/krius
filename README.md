# krius

<img src="./images/krius-logo.jpg" width="300" height="200">

Krius is a cli tool to setup Prometheus, Thanos &amp; friends across multiple clusters easily for scale

# Installation

Currently supported installation is source installation and please follow the steps as follows

## Compiling from source

#### Step 1: Clone the repository

```bash
$ git clone https://github.com/infracloudio/krius.git
```

#### Step 2: Build binary using make

```bash
$ make
```

Note: This will be installed in the respective path: $GOPATH/bin

## CLI Usage

```bash
$ krius --help
A tool to setup Prometheus, Thanos & friends across multiple clusters easily for scale .

Usage:
  krius [command]

Available Commands:
  configure   Configure the give component
  help        Help about any command
  install     Install the given component
  spec        Profile to be created
  uninstall   Deletes the installed stack

Flags:
  -h, --help   help for krius

Use "krius [command] --help" for more information about a command.
```

### Generate a spec file

Create a spec file based on the mode. This decides how to send data from a Prometheus instance to object storage one is using sidecar (pull model) and other is using receiver (push model). So the user has to choose one of two options based on needs and the limitations in underlying infrastructure.

```bash
$ krius spec generate --mode <receiver/sidecar>
```

### Spec File Configuration Details

```
---
clusters:
  - name: kind-cluster1
    type: prometheus
    data:
      install: true
      name: krius-prometheus
      namespace: monitoring
      mode: receiver
      objStoreConfig: krius-bucket
  - name: kind-cluster2
    type: thanos
    data:
      name: kind-thanos
      namespace: monitoring
      objStoreConfig: krius-bucket
      querier:
        name: global
        dedupEnbaled: true
        autoDownSample: true
        partialResponse: true
      querierFE:
        name: testing
        cacheOption: inMemory
      receiver:
        name: test
      compactor:
        name: test
        downsampling: true
        deduplication: true
        retentionResolutionRaw: 30d
        retentionResolution5m: 30d
        retentionResolution1h: 10y
      ruler:
        alertManagers:
          - http://kube-prometheus-alertmanager.monitoring.svc.cluster.local:9093
        config: |-
          groups:
            - name: "metamonitoring"
              rules:
                - alert: "PrometheusDown"
                  expr: absent(up{prometheus="monitoring/kube-prometheus"})
objStoreConfigslist:
  - name: krius-bucket
    type: s3
    config:
      bucket: name-of-bucket
      endpoint: s3.us-west-2.amazonaws.com
      accessKey: your-access-key-id
      secretKey: your-secret-access-key
      insecure: false
    bucketweb:
      enabled: false

```

- clusters: a list of kubernetes clusters to setup promethues or thanos
  - name: Kubernetes context name
  - type: Name of setup Prometheus or Thanos
  - data: Config for the setup type
- objStoreConfigslist - a list of object storage buckets. We refer them in cluster spec by name
  - name: Any unique name to refer this in cluster
  - type: Supported clients- AWS/S3 (and all S3-compatible storages e.g Minio)
  - config: Provide bucket, endpoint, accessKey, and secretKey keys to access storage

### Deploy Krius stack by applying the generated spec file

This command will validate the spec file passed as config-file and apply the configuration.

```bash
$ krius spec apply --config-file <relative-path/filename>
```

### Uninstall Krius stack using the spec file

This command will deletes the entire stack across clusters.

```bash
$ krius spec uninstall --config-file <relative-path/filename>
```

### Describe Krius Stack using the spec file[WIP]

This command will describe the clusters with meta-data details.

```bash
$ krius spec describe-cluster --config-file <relative-path/filename>
```