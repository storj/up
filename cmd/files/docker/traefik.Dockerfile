FROM traefik:2.4.8

COPY traefik/ /
COPY certificates /etc/traefik/certificates
