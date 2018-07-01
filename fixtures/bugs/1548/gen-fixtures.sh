#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
testcases="fixture-1548.yaml"
for testcase in ${testcases} ; do
    target=gen-${testcase%.yaml}
    spec=./${testcase}
    serverName="base64"
    logfile=${testcase%.yaml}.log
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --spec ${spec} --name=${serverName} --target ${target} --output=${logfile}
    if [[ $? != 0 ]] ; then
        echo "Server generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Server generation OK"
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
