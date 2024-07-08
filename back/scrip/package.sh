#!/bin/bash

cp scrip/application-web.service scrip/install.sh  scrip/uninstall.sh scrip/get_frame scrip/upgrade.sh config/application-web.yaml database/application-web.db .

mkdir -p application_web
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o application-web -trimpath -ldflags '-s -w'  
mv application-web application-web.yaml application-web.db application-web.service get_frame  install.sh upgrade.sh  uninstall.sh application_web
tar -czvf application-web-linux_arm64.tgz application_web
rm -rf application_web