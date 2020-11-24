#! /bin/bash
rm -rf case1 && mkdir case1
swagger generate server \
  --spec fixture-2362.yaml \
  --target case1 \
  --model-package restapi/models \
  --principal github.com/go-swagger/go-swagger/fixtures/bugs/2362/case1/restapi/models.Principal \
  --exclude-main \
  --name data-updater
printf "package models\n\ntype Principal struct{}" >> case1/restapi/models/principal.go

rm -rf case2 && mkdir case2
swagger generate server \
  --spec fixture-2362.yaml \
  --target case2 \
  --model-package restapi/models \
  --principal models.Principal \
  --exclude-main \
  --name data-updater
printf "package models\n\ntype Principal struct{}" >> case2/restapi/models/principal.go

rm -rf case3 && mkdir case3
swagger generate server \
  --spec fixture-2362.yaml \
  --target case3 \
  --model-package restapi/models \
  --principal restapi/models.Principal \
  --exclude-main \
  --name data-updater
printf "package models\n\ntype Principal struct{}" >> case3/restapi/models/principal.go

rm -rf case4 && mkdir case4
swagger generate server \
  --spec fixture-2362.yaml \
  --target case4 \
  --model-package restapi/models \
  --principal internal.Principal \
  --exclude-main \
  --name data-updater
mkdir -p case4/internal
printf "package internal\n\ntype Principal struct{}" >> case4/internal/principal.go

rm -rf case5 && mkdir case5
swagger generate server \
  --spec fixture-2362.yaml \
  --target case5 \
  --model-package restapi/generated \
  --principal models.Principal \
  --exclude-main \
  --name data-updater
mkdir -p case5/models
printf "package models\n\ntype Principal struct{}" >> case5/models/principal.go
