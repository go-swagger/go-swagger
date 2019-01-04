#! /bin/bash 
build_fixtures() {
_artefact=$1
echo "Building ${_artefact}(s) for fixtures..."
shift
_testcases=$*

for testcase in ${_testcases} ; do
    target=./gen-${testcase%.*}
    spec=${testcase}
    serverName="nrcodegen"
    logfile=${testcase%.*}.log 
    serverDir=${target}/cmd/${serverName}"-server"
    rm -rf ${target}
    mkdir ${target}
    echo "${_artefact} generation for ${spec}"
    swagger generate ${_artefact} --spec ${spec} --target ${target}  --name=${serverName} --output=${logfile}  
    #--skip-flatten
    if [[ $? != 0 ]] ; then
        echo "${_artefact} generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: ${_artefact} generation OK"
    if [[ ! -d ${serverDir} ]] ; then
        echo "No server found!"
        exit 1
    else
        (cd ${serverDir} ; go build)
        if [[ $? != 0 ]] ; then
            echo "${_artefact} build failed for ${spec}"
            exit 1
        fi
        echo "${spec}: ${_artefact} build OK"
        if [[ -n ${clean} ]] ; then 
             rm -rf ${target}
             rm -f ${logfile}
        fi
    fi
done
}
