<p align="center">  
  <img alt="krius" src="./images/krius-02.svg" width="250" height="250">
</p>

![Go Report Card](https://goreportcard.com/badge/github.com/infracloudio/krius)  ![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg) ![workflow](https://github.com/infracloudio/krius/actions/workflows/workflow.yaml/badge.svg)



# Krius

Krius is a cli tool to setup and manage Prometheus, Thanos &amp; friends across multiple clusters easily for scale

## Prerequisites

- Kubernetes 1.16+
- Helm 3+

## Installation

## Install the pre-compiled binary

```bash
curl -Lsf https://raw.githubusercontent.com/infracloudio/krius/main/install.sh | bash
```

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
      install: true
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
        config:
          maxSixe: 1
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
        name: ruler
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
    type: S3
    config:
      bucket: name-of-bucket
      endpoint: s3.us-west-2.amazonaws.com
      access_key: your-access-key-id
      secret_key: your-secret-access-key
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

![](./images/spec-apply.gif)

### Uninstall Krius stack using the spec file

This command will deletes the entire stack across clusters.

```bash
$ krius spec uninstall --config-file <relative-path/filename>
```

![](./images/spec-destroy.gif)

### Describe Krius Stack using the spec file[WIP]

This command will describe the clusters with meta-data details.

```bash
$ krius spec describe-cluster --config-file <relative-path/filename>
```

### Upgrade Krius stack using the spec file

1. To upgrade any cluster component, mark install to false, do the required change in cluster spec and run the spec apply command
2. To upgrade the object storage for a cluster, mark install to false, update the objStoreConfigslist's config, change the type and run the spec apply command

```
  data:
      install: false
```

```bash
$ krius spec apply --config-file <relative-path/filename>
```

for example if you want to upgrade the long-term object storage back-end from AWS S3 to GCP storage buckets

1. update the object storage config
2. mark cluster.data install to false
3. run spec apply command

Sample GCS config

```
type: GCS
config:
  bucket: "krius-bucket"
  service_account: |-
    {
      "type": "service_account",
      "project_id": "project",
      "private_key_id": "abcdefghijklmnopqrstuvwxyz12345678906666",
      "private_key": "-----BEGIN PRIVATE KEY-----\...\n-----END PRIVATE KEY-----\n",
      "client_email": "project@thanos.iam.gserviceaccount.com",
      "client_id": "123456789012345678901",
      "auth_uri": "https://accounts.google.com/o/oauth2/auth",
      "token_uri": "https://oauth2.googleapis.com/token",
      "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
      "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/thanos%40gitpods.iam.gserviceaccount.com"
    }

```

## Contributing

Contributions are welcome, please [read the guidelines](CONTRIBUTING.md)
