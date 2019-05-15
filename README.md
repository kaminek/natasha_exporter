# Natasha Exporter

A [Prometheus](https://prometheus.io/) exporter for
[Natasha](https://github.com/scaleway/natasha) that collects metrics

## Metrics

This exporter expose the following Natasha metrics:
* Metrics about [DPDK](htpps://dpdk.org) port statistics (stats and xstats)
* Metrics about [Natasha](https://github.com/scaleway/natasha) application stats.
* Others

For more information about [Natasha](https://github.com/scaleway/natasha)

## Install

You can download prebuilt binaries from our [GitHub
releases](https://github.com/kaminek/natasha_exporter/releases). Or build it
yourself using
```
make install
```

## Development

Make sure you have a working Go environment, for further reference or a guide
take a look at the [install instructions](http://golang.org/doc/install.html).
This project requires Go >= v1.8.

```bash
go get -d github.com/kaminek/natasha_exporter
cd $GOPATH/src/github.com/kaminek/natasha_exporter

# get deps
make dep

# build binary
make

./bin/natasha_exporter -h
```

## Authors

* [Amine Kherbouche](https://github.com/kaminek)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2019 Amine KHERBOUCHE <akherbouche@scaleway.com>
```
