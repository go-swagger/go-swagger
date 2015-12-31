#!/bin/sh


apt-get update -yqq
apt-get install jq unzip

os=$(uname|tr '[:upper:]' '[:lower:]')
latestv=$(curl -s https://api.github.com/repos/Masterminds/glide/releases/latest | jq -r .tag_name)
curl -o /tmp/glide-linux.zip -sL "https://github.com/Masterminds/glide/releases/download/$latestv/glide-$latestv-$os-amd64.zip"
unzip -j -d /usr/bin /tmp/glide-linux.zip
