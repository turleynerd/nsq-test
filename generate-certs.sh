#!/bin/bash

# create the root key
openssl genrsa -aes256 -out certs/root.key 4096

# create and self sign the root certificate
openssl req -x509 -new -nodes -key certs/root.key -sha256 -days 1024 -out certs/root.crt

# create the certiicate key
openssl genrsa -out certs/make-production.internal.key 2048

# create the signing request
openssl req -new -sha256 -key certs/make-production.internal.key -subj "/C=US/ST=NC/O=Red Ventures/CN=nsqd" -out certs/make-production.internal.csr

# create the certificate using the csr and certificate key along with the root key
openssl x509 -req -in certs/make-production.internal.csr -CA certs/root.crt -CAkey certs/root.key -CAcreateserial -out certs/make-production.internal.crt -days 500 -sha256
