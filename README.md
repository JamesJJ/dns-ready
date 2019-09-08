# DNS Ready

[![Docker Automated build](https://img.shields.io/docker/cloud/automated/jamesjj/dns-ready)](https://hub.docker.com/r/jamesjj/dns-ready/)
[![Docker Automated build](https://img.shields.io/docker/cloud/build/jamesjj/dns-ready)](https://hub.docker.com/r/jamesjj/dns-ready/)


Sometimes we just need to wait until a hostname is resolvable in DNS.

  * This will repeatedly attempt to resolve a hostname and then exit gracefully if the DNS lookup is successful, or when the maximum number of attempts is reached.

  * TERM and INT signals will cause immediate graceful exit.

  * Graceful exit means ending the program with return code zero.

  * This can be used as an "init container" in Kubernetes ([docker image](https://hub.docker.com/r/jamesjj/dns-ready): `jamesjj/dns-ready`).


## Configuration:

*Options can be configured using command line flags, or `DNSREADY_*` environment variables*:

| Flag           | Environment variable   | Default                                | Description                                         |
|----------------|------------------------|----------------------------------------|-----------------------------------------------------|
| `-acceptempty` | `DNSREADY_ACCEPTEMPTY` | false                                  | Accept a DNS response with no IP addresses          |
| `-host`        | `DNSREADY_HOST`        | kube-dns.kube-system.svc.cluster.local | The hostname to resolve                             |
| `-pause`       | `DNSREADY_PAUSE`       | 800                                    | Milliseconds to sleep between attempts              |
| `-timeout`     | `DNSREADY_TIMEOUT`     | 1200                                   | Timeout in milliseconds for each DNS lookup attempt |
| `-retries`     | `DNSREADY_RETRIES`     | 30                                     | Maximum number of attempts before graceful exit     |
| `-verbose`     | `DNSREADY_VERBOSE`     | false                                  | Show each attempt on STDOUT                         |
| `-silent`      | `DNSREADY_SILENT`      | false                                  | Do not show anything on STDOUT                      |

