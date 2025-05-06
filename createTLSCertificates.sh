#!/bin/sh

# Usage: ./createTLSCertificates.sh <SERVER_NAME> [<CA_NAME>]
# Example: ./createTLSCertificates.sh myserver myca

mkdir certs &> /dev/null
cd certs
SERVER_NAME="${1:-server}"
CA_NAME="${2:-myca}"

# Verzeichnis erstellen
mkdir -p "${CA_NAME}_certs"
cd "${CA_NAME}_certs" || exit 1

# CA erstellen (falls nicht vorhanden)
if [ ! -f "${CA_NAME}.key" ]; then
  echo "Generating CA key and certificate..."
  openssl ecparam -genkey -name prime256v1 -out "${CA_NAME}.key"
  openssl req -x509 -new -key "${CA_NAME}.key" -out "${CA_NAME}.crt" \
    -subj "/CN=${CA_NAME}_CA/O=${CA_NAME}_Org" -days 3650
fi

# Server-Schlüssel und CSR
echo "Generating server key and CSR..."
openssl ecparam -genkey -name prime256v1 -out "${SERVER_NAME}.key"
openssl req -new -key "${SERVER_NAME}.key" -out "${SERVER_NAME}.csr" \
  -subj "/CN=${SERVER_NAME}/O=${SERVER_NAME}_Org"

# SAN-Konfiguration (für localhost/IP)
cat > "${SERVER_NAME}.cnf" <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req

[v3_req]
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${SERVER_NAME}
DNS.2 = localhost
IP.1 = 127.0.0.1
EOF

# Zertifikat signieren
echo "Signing server certificate..."
openssl x509 -req -in "${SERVER_NAME}.csr" \
  -CA "${CA_NAME}.crt" -CAkey "${CA_NAME}.key" -CAcreateserial \
  -out "${SERVER_NAME}.crt" -days 365 \
  -extfile "${SERVER_NAME}.cnf" -extensions v3_req

echo "✅ Done! Server files:"
ls -l "${SERVER_NAME}".{key,crt}
