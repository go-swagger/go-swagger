#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
testcases="fixture-910.yaml fixture-910-2.yaml"
for testcase in ${testcases} ; do
    target=./gen-${testcase%.yaml}
    spec=./${testcase}
    serverName="nrcodegen"
    logfile=${testcase%.yaml}".log"
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --spec ${spec} --target ${target} --name ${serverName} --output=${logfile}
    if [[ $? != 0 ]] ; then
        echo "Generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Generation OK"
    (cd ${target}/cmd/${serverName}"-server"; go build)
    if [[ $? != 0 ]] ; then
        echo "Build failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Build OK"
    if [[ -n ${clean} ]] ; then 
        rm -rf ${target}
    fi
done
