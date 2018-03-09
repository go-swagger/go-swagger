#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
testcases="fixture-1314.yaml"
for testcase in ${testcases} ; do
    target=./gen-${testcase%.yaml}
    spec=./${testcase}
    serverName="nrcodegensrv"
    srvdir=${target}/cmd/${serverName}-server
    logfile=${testcase%.yaml}.log
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --spec ${spec} --target ${target} --output ${logfile} --name ${serverName}
    if [[ $? != 0 ]] ; then
        echo "Generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Generation OK"
    if [[ -d ${srvdir}  ]]; then
      (cd ${srvdir}; go build)
      if [[ $? != 0 ]] ; then
        echo "Build failed for ${spec}"
        exit 1
      fi
      echo "${spec}: Build OK"
      if [[ -n ${clean} ]] ; then 
        rm -rf ${target}
      fi
    else
      echo "No server generated for ${spec}"
      exit 1
    fi
done
