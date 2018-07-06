# coredns-register

Coredns-register is an agent that registers records in CoreDNS with etcd as the backend.
Provides dynamic service registration when using CoreDNS for DNS discovery of prometheus.

## How to use

```shell
coredns-register -config /etc/coredns-register/config.yml
```

### Using Docker

```shell
docker run -d \
  --env HOSTNAME=agent-host \
  --env ADDRESS=127.0.0.1 \
  --env DISCOVERY_SRV=srv.coredns.test \
  --env SRV_RECORDS="node.coredns.test:9100,etcd.coredns.test:2379"
  kanga333/coredns-register
```

Refer to config.yml for configurable environment variables.

## Settings

yaml example as below.

```yaml:/etc/coredns-register/config.yml
hostname: host # The hostname of the FQDN.
address: 127.0.0.1 # IP address resolved by DNS.
interval: 60 # Registration interval.
etcd:
  discovery-srv: dns.domain.test # SRV indicating the destination of etcd used by CoreDNS.
  endpoints: # Directly specify the destination of etcd used by CoreDNS. discovery-srv takes precedence.
    - http://127.0.0.1:2379
  basepath: /skydns # Base path of etcd registration.
record_files:
  - "/etc/coredns-register/record.d/*yml" # Path of srv record configuration file.
records:
  srv:
    - domain: a.domain.test # The domain of the FQDN.
      address: 127.0.0.1 # Override setting root adress. Option.
      port: 80 # Port of SRV record.
srv_records: "c.domain.test:80,d.domain.test:80" # Comma separated list of "domain:port" of srv record. Option.
```

```yaml:/etc/coredns-register/record.d/b.yml
srv:
  - domain: b.domain.test # The domain of the FQDN.
    address: 127.0.0.1 # Override setting root adress. Option.
    port: 80 # Port of SRV record.
```
