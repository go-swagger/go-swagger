#! /bin/bash
#swagger.exe generate server -f swagger.yml /A vg-api /m restapimodels
rm -rf tmp && mkdir tmp
swagger generate server -f fixture-1808.yaml -A vg-api -m restapimodels --target tmp
