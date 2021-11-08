#!/bin/bash
# Generate a certificate for the authservice and traefik proxy.
# These are hostnames inside the docker network.
# Future improvements:
# - Run a DNS server for a tld like *.storj.test so it can be accessed from the outside.
# - Have a single CA certificate to sign multiple certificates for all the services.
openssl req \
    -x509 \
    -newkey rsa:4096 \
    -keyout up-authservice-key.pem \
    -out up-authservice-cert.pem \
    -days 3650 \
    -nodes \
    -subj '/CN=localhost' \
    -addext "subjectAltName = DNS:localhost,DNS:traefik-authservice:DNS:authservice"
