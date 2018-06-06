#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
testcases="fixture-1400.json"
for testcase in ${testcases} ; do
    target=gen-${testcase%.json}
    spec=./${testcase}
    serverName="nrcodegen"
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --spec ${spec} --target ${target} --name=${serverName} --quiet
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
# test query
./${target}/cmd/${serverName}"-server"/${serverName}"-server" --port=8080 --scheme=http &
pid=$!
sleep 2
curl -X POST -v -H "Content-Type: multipart/form-data; boundary=------------------------a31e2ddd4b2c0d92" http://localhost:8080/file
kill ${pid}
