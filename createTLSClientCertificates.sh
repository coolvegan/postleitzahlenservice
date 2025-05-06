#!/bin/bash

# Usage: ./createTLSClientCertificates.sh <CLIENT_NAME> [<CA_NAME>]
# Example: ./createTLSClientCertificates.sh myca

mkdir certs &> /dev/null
cd certs

CLIENT_NAME="${1:-client}"
CA_NAME="${2:-myca}"

# Verzeichnis prüfen
if [ ! -d "${CA_NAME}_certs" ]; then
  echo "❌ Error: CA directory '${CA_NAME}_certs' not found. Run createTLSCertificates.sh first!"
  exit 1
fi

cd "${CA_NAME}_certs" || exit 1

# Client-Schlüssel und CSR
echo "Generating client key and CSR..."
openssl ecparam -genkey -name prime256v1 -out "${CLIENT_NAME}.key"
openssl req -new -key "${CLIENT_NAME}.key" -out "${CLIENT_NAME}.csr" \
  -subj "/CN=${CLIENT_NAME}/O=${CLIENT_NAME}_Org"

# Client-Zertifikat signieren (mit clientAuth)
echo "Signing client certificate..."
openssl x509 -req -in "${CLIENT_NAME}.csr" \
  -CA "${CA_NAME}.crt" -CAkey "${CA_NAME}.key" -CAcreateserial \
  -out "${CLIENT_NAME}.crt" -days 365 \
  -extensions clientAuth

echo "✅ Done! Client files:"
ls -l "${CLIENT_NAME}".{key,crt}
