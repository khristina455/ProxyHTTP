#!/bin/sh

openssl req -new -newkey rsa:2048 -nodes -keyout /certs/private_key.key -out /certs/CSR.csr -sha256 -subj "/CN=$1"
openssl x509 -req -days 365 -in /certs/CSR.csr -CA /certs/CA_certificate.crt -CAkey /certs/CA_private_key.key -out /certs/certificate.crt -set_serial $2 -sha256
