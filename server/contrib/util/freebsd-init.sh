#!/bin/sh

set -e # Exit script on error

sudo portsnap fetch extract

cd /usr/ports/www/gohugo && \
sudo make install clean BATCH=yes

# At this point, DNS records must be configured in order to use LetsEncrypt for deployment.
echo -e "\n\033[1;35mDNS records must be configured by now.\033[0m"
