# prom-timestream

This is an adapter for [Prometheus](https://prometheus.io) to store data in [Amazon Timestream](https://aws.amazon.com/timestream/). It currently supports the remote write API with the remote read being added soon.

## Build

```
$ make
```

## Run

The server supports the following options:
```
Usage of bin/server:
  -addr string
    	 (default ":4000")
  -database-name string
    	The database name to use in timestream (default "prom")
  -table-name string
    	The table name to use in timestream (default "metrics")
```

The following assumes you have a local AWS config ready with a profile `dev` and a timestream database `prom` with a table `metrics` already set up.
Below is an example on how to start the adapter on port 4000 on localhost.

```
$ AWS_SDK_LOAD_CONFIG=true AWS_REGION=eu-central-1 AWS_PROFILE=dev bin/server -database-name=prom -table-name=metrics
```

Futhermore the [test](/test) contains an example using docker-compose. In [prom.yml](/test/prom.yml) you find an example on how to configure prometheus for this adapter.
