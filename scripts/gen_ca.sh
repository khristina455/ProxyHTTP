#!/bin/sh

openssl req -new -newkey rsa:2048 -nodes -out CA_CSR.csr -keyout CA_private_key.key -subj "/CN=proxy"
openssl x509 -signkey CA_private_key.key -days 365 -req -in CA_CSR.csr -out CA_certificate.crt
cp CA_certificate.crt /usr/local/share/ca-certificates
update-ca-certificates
