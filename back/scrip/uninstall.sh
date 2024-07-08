#!/bin/sh

systemctl stop application-web
systemctl disable application-web
rm -rf /etc/application-web /var/lib/application-web /var/log/application-web /etc/systemd/system/application-web.service  /bin/application-web