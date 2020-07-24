#! /bin/bash
rm -rf internal && mkdir internal
swagger generate server \
  --spec fixture-2362.yaml \
  --target internal \
  --model-package restapi/models \
  --principal internal/restapi/models.Principal \
  --exclude-main \
  --name data-updater
