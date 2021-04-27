# krius
A tool to setup Prometheus, Thanos &amp; friends across multiple clusters easily for scale

## Prometheus Installation

#### Step 1: Clone the repo

```bash
$ git clone https://github.com/infracloudio/krius.git
```

#### Step 2: Build binary using make

```bash
$ make
```

#### Step 3: Run the command
```bash
$ krius install prometheus prom --namespace=demo
```