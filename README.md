# prom-timestream

A remote read/write adapter for using AWS Timestream as storage backend for
Prometheus.

## Build

```
$ make
```

## Run

```
$ AWS_SDK_LOAD_CONFIG=true AWS_REGION=eu-central-1 AWS_PROFILE=dev bin/server
```
