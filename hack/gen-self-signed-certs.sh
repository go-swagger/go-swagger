#!/bin/bash

go install github.com/cloudflare/cfssl/cmd/...@latest

csr=$(cat <<EOF
{
    "hosts": [
        "goswagger.local",
        "www.example.com",
        "https://www.example.com",
        "localhost",
        "127.0.0.1"
    ],
    "key": {
        "algo": "rsa",
        "size": 4096
    },
    "names": [
        {
            "C":  "US",
            "L":  "San Francisco",
            "O":  "go-swagger",
            "OU": "go-swagger",
            "ST": "California"
        }
    ]
}
EOF
)
TEMP=/tmp/certs
mkdir "${TEMP}"
CSR="${TEMP}/csr.json"
echo "${csr}" > "${CSR}"
cat "${CSR}"
pushd  "${TEMP}"
cfssl genkey -initca "${CSR}" | cfssljson -bare ca
cfssl gencert -ca "${TEMP}"/ca.pem -ca-key "${TEMP}"/ca-key.pem "${CSR}" |cfssljson -bare server
cfssl gencert -ca "${TEMP}"/ca.pem -ca-key "${TEMP}"/ca-key.pem "${CSR}" |cfssljson -bare client
popd

# restore keys and certs with their expected name
cp "${TEMP}"/ca-key.pem myCA.key
cp "${TEMP}"/ca.pem myCA.crt
cp "${TEMP}"/server-key.pem mycert1.key
cp "${TEMP}"/server.pem mycert1.crt
cp "${TEMP}"/client-key.pem myclient.key
cp "${TEMP}"/client.pem myclient.crt

rm -rf "${TEMP}"
