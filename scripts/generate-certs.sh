#!/bin/sh

CN="Wise2c CA"
# generate ca.key
openssl genrsa -out ca.key 2048
# generate ca.crt => caBundle
openssl req -x509 -new -nodes -key ca.key -subj "/CN=${CN}" -days 36500 -out ca.crt
# generate server.key => tls.key
openssl genrsa -out server.key 2048
# generate csr.conf
cat << EOF > csr.conf
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C = CN
ST = ShenZhen
L = SZ
O = Wise2c
OU = Wise2c
CN = ${CN}

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = admission-webhook-demo
DNS.2 = admission-webhook-demo.default
DNS.3 = admission-webhook-demo.default.svc
IP.1 = 127.0.0.1

[ v3_ext ]
authorityKeyIdentifier=keyid,issuer:always
basicConstraints=CA:FALSE
keyUsage=keyEncipherment,dataEncipherment
extendedKeyUsage=serverAuth,clientAuth
subjectAltName=@alt_names
EOF

# generate server.csr
openssl req -new -key server.key -out server.csr -config csr.conf
# generate server.crt => tls.crt
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 36500 -extensions v3_ext -extfile csr.conf
