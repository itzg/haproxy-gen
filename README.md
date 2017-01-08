This tool uses a Go template and expands out the contents
for any number of domains to proxy.

Run it without any arguments to view usage:

```
Usage:
  haproxy-gen [command]

Available Commands:
  certbot-args   Generate the -d arguments suitable for passing to certbot
  generate       Generates the haproxy.cfg file
  primary-domain Extract just the primary domain from the given configuration and domains

Flags:
      --config string   config file (default is $HOME/.haproxy-gen.yaml)
  -t, --toggle          Help message for toggle

Use "haproxy-gen [command] --help" for more information about a command.
```

## Example

With the appropriate binary downloaded from the [releases](https://github.com/itzg/haproxy-gen/releases) and the [template file downloaded](https://raw.githubusercontent.com/itzg/haproxy-gen/master/haproxy.cfg.tmpl), running:

```
➜ chmod +x haproxy-gen_darwin_amd64
➜ ./haproxy-gen_darwin_amd64 generate -d testing.example.com@apache:80
```

produces the output:
```
global
    # set default parameters to the modern configuration
    tune.ssl.default-dh-param 2048
    ssl-default-bind-ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHAC
    ssl-default-bind-options no-sslv3 no-tlsv10 no-tlsv11 no-tls-tickets
    ssl-default-server-ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CH
    ssl-default-server-options no-sslv3 no-tlsv10 no-tlsv11 no-tls-tickets


defaults
    mode    http
    option  forwardfor
    option  http-server-close

    stats   enable
    stats   uri /hastats
    stats   realm HAproxy\ Stats
    stats   auth admin:haproxy

    timeout client 30s
    timeout connect 4s
    timeout server 30s

frontend ft
    bind    :443 ssl crt /etc/certs
    bind    :80
    redirect scheme https code 301 if !{ ssl_fc }
    # HSTS (15768000 seconds = 6 months)
    http-response set-header Strict-Transport-Security max-age=15768000


    acl host_testing_example_com hdr(host) -i testing.example.com
    use_backend testing_example_com if host_testing_example_com


backend testing_example_com
    server testing_example_com apache:80
```
