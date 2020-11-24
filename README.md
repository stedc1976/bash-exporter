# Prometheus bash exporter

Simple & minimalistic Prometheus exporter for bash scripts.

[![Go Report Card](https://goreportcard.com/badge/github.com/stedc1976/bash-exporter)](https://goreportcard.com/badge/github.com/stedc1976/bash-exporter)

## Installation

```console
$ docker build --rm -t diclem27/bash-exporter:1.0.0 .
```

## Docker quick start

```console
$ docker run -d -p 9300:9300 --name my_bash-exporter diclem27/bash-exporter:1.0.0
```

```console
$ curl -s 127.0.0.1:9300/metrics | grep ^bash
bash{container_name="sti-build",job="job1",namespace="service-activator",pod_name="samigrationpam-kieserver-26-build",verb:"row_count"} 3596
bash{container_name="sti-build",job="job1",namespace="service-activator",pod_name="samigrationpam-kieserver-26-build",verb:"row_count"} 3721
bash{container_name="sti-build",job="job1",namespace="service-activator",pod_name="samigrationpam-kieserver-26-build",verb:"row_count"} 3802
```

## Usage

```console
Usage of ./bash-exporter:
  -debug
    	Debug log level
  -interval int
    	Interval for metrics collection in seconds (default 5)
  -path string
    	path to directory with bash scripts (default "/scripts")
  -web.listen-address string
    	Address on which to expose metrics (default ":9300")
```

Just point `-path` flag to the directory with your bash scripts. Names of the files (`(.*).sh`) will be used as the `job` label. Bash scripts should return valid json.

## External doc

https://godoc.org/github.com/prometheus/client_golang/prometheus

## TODO
- [] Helm Chart
- [] Several scripts
