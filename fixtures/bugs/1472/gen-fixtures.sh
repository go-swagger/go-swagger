#! /bin/bash 
rm -rf ./restapi/

swagger generate server --exclude-main -m restapi/models -f bug_1472.yml
