#! /bin/bash 
curl \
  --verbose \
  --get \
  http://localhost:43016/api/items

basic=`echo "fred:scrum"|tr -d '\n'|base64 -i`
curl \
  --verbose \
  --get \
  --header "Authorization: Basic ${basic}" \
  http://localhost:43016/api/account

basic=`echo "fred:scrum"|tr -d '\n'|base64 -i`
curl \
  --verbose \
  --get \
  --header "Authorization: Basic ${basic}" \
  "http://localhost:43016/api/order/myOrder?access_token=`cat tokens/token-apikey-customer.jwt`"

basic=`echo "ivan:terrible"|tr -d '\n'|base64 -i`
curl \
  --verbose \
  -i \
  -X POST \
  --data '{"orderID": "myorder", "orderLines": [{"quantity": 10, "purchasedItem": "myItem"}]}' \
  --header "Content-Type: application/json" \
  --header "Authorization: Basic ${basic}" \
  "http://localhost:43016/api/order/add?access_token=`cat tokens/token-apikey-customer.jwt`"

basic=`echo "ivan:terrible"|tr -d '\n'|base64 -i`
curl \
  --verbose \
  -i \
  -X POST \
  --data '{"orderID": "myorder", "orderLines": [{"quantity": 10, "purchasedItem": "myItem"}]}' \
  --header "Content-Type: application/json" \
  --header "Authorization: Basic ${basic}" \
  --header "X-Custom-Key: `cat tokens/token-apikey-reseller.jwt`" \
  "http://localhost:43016/api/order/add?access_token=`cat tokens/token-bearer-inventory-manager.jwt`"

curl \
  --verbose \
  --get \
  --header "X-Custom-Key: `cat tokens/token-apikey-reseller.jwt`" \
  "http://localhost:43016/api/orders/myItem"

