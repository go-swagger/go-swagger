#! /bin/bash 

# jti: ID
# iss: Issuer
# roles: custom claim
#
# Token for inventoryManager
token='token-bearer-inventory-manager.jwt'
echo \
'{"jti": "fred", "iss": "example.com", "roles": [ "inventoryManager" ]}'|\
jwt -key ../keys/apiKey.prv -alg RS256 -sign - > ${token}
jwt -key ../keys/apiKey.pem -alg RS256 -verify ${token}
jwt -show ${token}

# Token for API Keys
token='token-apikey-reseller.jwt'
echo \
'{"jti": "fred", "iss": "example.com", "roles": [ "reseller" ]}'|\
jwt -key ../keys/apiKey.prv -alg RS256 -sign - > ${token}
jwt -key ../keys/apiKey.pem -alg RS256 -verify ${token}
jwt -show ${token}

token='token-apikey-customer.jwt'
echo \
'{"jti": "ivan", "iss": "example.com", "roles": [ "customer" ]}'|\
jwt -key ../keys/apiKey.prv -alg RS256 -sign - > ${token}
jwt -key ../keys/apiKey.pem -alg RS256 -verify ${token}
jwt -show ${token}
