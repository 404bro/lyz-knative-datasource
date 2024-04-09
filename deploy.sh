#!/bin/bash

TARGET="node-t"
TARGET_ZIP="/home/lyz/lyz-knative-datasource.zip"
TARGET_PLUGINS_DIR="/home/lyz/plugins"

mage -v
npm run build
mv dist/ lyz-knative-datasource
zip -0 lyz-knative-datasource.zip lyz-knative-datasource -r
scp lyz-knative-datasource.zip $TARGET:$TARGET_ZIP

ssh $TARGET "rm -rf $TARGET_PLUGINS_DIR/lyz-knative-datasource"
ssh $TARGET "unzip $TARGET_ZIP -d $TARGET_PLUGINS_DIR/lyz-knative-datasource"
ssh $TARGET "chmod -R 777 $TARGET_PLUGINS_DIR/lyz-knative-datasource"
ssh $TARGET "kubectl delete pod -n default -l app.kubernetes.io/name=grafana"
