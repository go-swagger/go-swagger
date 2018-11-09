#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
testcases="fixture-1624.json"
for testcase in ${testcases} ; do
    target=gen-${testcase%.json}
    spec=./${testcase}
    serverName="nrcodegen"
    serverPackage="pkg/to/nrcodegen"
    configureFile="${serverPackage}/configure_nrcodegen.go"
    logfile=${testcase%.json}.log
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --spec ${spec} --name=${serverName} --target ${target} --output=${logfile} --server-package=${serverPackage}
    if [[ $? != 0 ]] ; then
        echo "Server generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Server generation OK"
    
    grep -Fq "${serverPackage}" ${target}/cmd/${serverName}"-server"/main.go
    if [[ $? != 0 ]] ; then
        echo "${spec}: imports using ServerPackage failed for ${spec}"
        exit 1
    fi
    echo "${spec}: imports using ServerPackage OK"
    
    grep -Fxq "package ${serverName}" ${target}/${serverPackage}/configure_nrcodegen.go
    if [[ $? != 0 ]] ; then
        echo "${spec}: ServerPackage name failed for ${spec}"
        exit 1
    fi
    echo "${spec}: ServerPackage name OK"

    (cd ${target}/cmd/${serverName}"-server"; go build)
    if [[ $? != 0 ]] ; then
        echo "Server build failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Server build OK"
    if [[ -n ${clean} ]] ; then 
        rm -rf ${target}
    fi
done
