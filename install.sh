#!/bin/sh
set -e

go build -o gitwit .
sudo install gitwit /usr/local/bin/gitwit
sudo codesign -f -s - /usr/local/bin/gitwit
echo "gitwit installed to /usr/local/bin/gitwit"
