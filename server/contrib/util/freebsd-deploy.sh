#!/bin/sh

set -e # Exit script on error

# Destination paths.
: ${prefix="/usr/local"}
: ${dest_bin="${prefix}/bin"}
: ${dest_rc_d="${prefix}/etc/rc.d/bigtechfaild"}
: ${dest_rootdir="/var/www/bigtechfaild"}

# Project resources.
: ${bigtechfaild_bin="./server/bigtechfaild"}
: ${rc_d_script="./server/contrib/init/bigtechfaild.rc.d"}

# Build the server
echo -e "\033[1;35mBuilding and testing the server ... \033[0m"
cd server
go build
go test
cd ..

# Build the website
echo -e "\033[1;35mBuilding the website ... \033[0m"
cd website
hugo --minify
cd ..

# Install the database, program binary, and static files.
echo -e "\033[1;35mInstalling everything ... \033[0m"

sudo service bigtechfaild stop || echo

sudo cp ${bigtechfaild_bin} ${dest_bin}

sudo mkdir -p ${dest_rootdir}
sudo rm -rf ${dest_rootdir}/*
sudo cp -r ./website/public/* ${dest_rootdir}

# The rc.d script
sudo cp ${rc_d_script} ${dest_rc_d}
sudo sysrc bigtechfaild_enable=YES

sudo service bigtechfaild start
