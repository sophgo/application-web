#!/bin/bash

set -e
if [[ -f /etc/systemd/system/application-web.service ]]; then
  systemctl stop application-web.service
  systemctl disable application-web.service
fi
mkdir -p /etc/application-web/config /var/log/application-web /var/lib/application-web/db 

cp application-web /bin
cp application-web.yaml /etc/application-web/config
cp application-web.db /var/lib/application-web/db
cp application-web.service /etc/systemd/system/
cp get_frame /var/lib/application-web

systemctl daemon-reload
systemctl enable application-web.service
systemctl start application-web.service