#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
# Fixtures with models only
testcases="fixture-1392-1.yaml fixture-1392-2.yaml fixture-1392-3.yaml fixture-1392-4.yaml fixture-1392-5.yaml"
for testcase in ${testcases} ; do
    target=./gen-${testcase%.*}
    spec=./${testcase}
    serverName="codegensrv"
    rm -rf ${target}
    mkdir ${target}
    echo "Generation for ${spec}"
    swagger generate server --spec ${spec} --name=${serverName} --target ${target} --output=${testcase%.*}.log
    # 1>x.log 2>&1
    #
    if [[ $? != 0 ]] ; then
        echo "Generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Generation OK"
    if [[ ! -d ${target}/models ]] ; then
        echo "No model in this spec! Skipped"
    else
        echo "Building model package"
        (cd ${target}/models; go build)
        if [[ $? != 0 ]] ; then
            echo "Model build failed for ${spec}"
            failed=1
        fi

        echo "Build server"
        (cd ${target}/cmd/${serverName}"-server"; go build)
        if [[ $? != 0 ]] ; then
           echo "Server build failed for ${spec}"
           failed=1
        fi
        if [[ ${failed} == 1 ]] ; then 
            echo "Continue..."
            #exit 1
        else
        echo "${spec}: Build OK"
        if [[ -n ${clean} ]] ; then 
             rm -rf ${target}
        fi
        fi
    fi
done
# Non reg codegen
# NOTE(fredbi): 
# - azure: invalid spec 
# - bitbucket: model does not compile
# - issue72: invalid spec
# - todolist.discriminator: known issue with schemavalidator
testcases=`cd ../../codegen;ls -1|grep -vE 'azure|bitbucket|existing-model|issue72|todolist.simple.yml'`
for testcase in ${testcases} ; do
    target=./gen-${testcase%.*}
    if [[ -f ../../codegen/${testcase} ]] ; then
      spec=../../codegen/${testcase}
    else 
      spec=${testcase}
    fi
    serverName="nrcodegen"
    rm -rf ${target}
    mkdir ${target}
    echo "Generation for ${spec}"
    swagger generate server --skip-validation --spec ${spec} --target ${target} --name=${serverName} --output=${testcase%.*}.log
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
# 
testcases=
for testcase in ${testcases} ; do
    target=./gen-${testcase%.*}
    spec=./${testcase}
    serverName="bugfix"
    rm -rf ${target}
    mkdir ${target}
    echo "Generation for ${spec}"
    swagger generate server --spec ${spec} --target ${target} --quiet --name=${serverName}
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
