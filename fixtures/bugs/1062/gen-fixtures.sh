#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
testcases="eve-online-esi.json"
for testcase in ${testcases} ; do
    target=./gen-${testcase%.json}
    spec=./${testcase}
    serverName="nrcodegensrv"
    srvdir=${target}/cmd/${serverName}-server
    logfile=${testcase%.json}.log
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --spec ${spec} --target ${target} --output ${logfile} --name ${serverName}
    if [[ $? != 0 ]] ; then
        echo "Server generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Server generation OK"
    if [[ -d ${srvdir}  ]]; then
      (cd ${srvdir}; go build)
      if [[ $? != 0 ]] ; then
        echo "Server rbBuild failed for ${spec}"
        exit 1
      fi
      echo "${spec}: Server build OK"
      if [[ -n ${clean} ]] ; then 
        rm -rf ${target}
      fi
    else
      echo "No server generated for ${spec}"
      exit 1
    fi
    target=${target}"-client"
    logfile=${testcase%.json}"-client.log"
    rm -rf ${target}
    mkdir ${target}
    swagger generate client --spec ${spec} --target ${target} --output ${logfile} --name ${serverName}
    if [[ $? != 0 ]] ; then
        echo "Client generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Client generation OK"
    if [[ -d ${srvdir}  ]]; then
      (cd ${srvdir}; go build)
      if [[ $? != 0 ]] ; then
        echo "Client build failed for ${spec}"
        exit 1
      fi
      echo "${spec}: Client build OK"
      if [[ -n ${clean} ]] ; then 
        rm -rf ${target}
      fi
    else
      echo "No client generated for ${spec}"
      exit 1
    fi
done
