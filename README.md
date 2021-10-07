# krius

A tool to setup Prometheus, Thanos &amp; friends across multiple clusters easily for scale

## Compiling from source

#### Step 1: Clone the repo

```bash
$ git clone https://github.com/infracloudio/krius.git
```

#### Step 2: Build binary using make

```bash
$ make
```

## CLI Usage

```bash
$ krius --help
```

### Generate a spec file

Create a spec file based on the mode. This decides how to send data from a Prometheus instance to object storage one is using sidecar (pull model) and other is using receiver (push model). So the user has to choose one of two options based on needs and the limitations in underlying infrastructure.

```bash
$ krius spec generate --mode <receiver/sidecar>
```

### Lets go through the structure of the spec file

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
  - name: kubernetes context name
  - type: name of setup Prometheus or Thanos
  - data: config for the setup type
- objStoreConfigslist - a list of object storage buckets. We refer them in cluster spec by name
  - name: any unique name to refer this in cluster
  - type: supported clients- AWS/S3 (and all S3-compatible storages e.g Minio)
  - config: provide bucket, endpoint, accessKey, and secretKey keys to access storage

### Run the generated spec file

```bash
$ krius spec apply --config-file <filename>
```
