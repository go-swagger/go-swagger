#! /bin/bash
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
continueOnError=
# A small utility to build fixture servers
# Fixtures with models only
testcases="fixture-2330.yaml"

# Generation options
fullFlatten="--with-flatten=full"
withExpand="--with-expand"
minimal="--with-flatten=minimal"
for opts in ${fullFlatten} ${withExpand} ${minimal} ; do
for testcase in ${testcases} ; do
    spec=${testcase}
    testcase=`basename ${testcase}`
    if [[ -z ${opts} || ${opts} == ${minimal} ]] ; then
        target=./gen-${testcase%.*}-minimal
    elif [[ ${opts} ==  ${fullFlatten} ]] ; then
        target=./gen-${testcase%.*}-flatten
    else
        target=./gen-${testcase%.*}-expand
    fi
    serverName="codegensrv"
    rm -rf ${target}
    mkdir ${target}
    echo "Server generation for ${spec} with opts=${opts}"
    swagger generate server --skip-validation ${opts} --spec=${spec} --target=${target} --log-output=${testcase%.*}.log
    # 1>x.log 2>&1
    #
    if [[ $? != 0 ]] ; then
        echo "Generation failed for ${spec}"
        if [[ ! -z ${continueOnError} ]] ; then
            failures="${failures} codegen:${spec}"
            continue
        else
            exit 1
        fi
    fi
    echo "${spec}: Generation OK"
    if [[ ! -d ${target}/models ]] ; then
        echo "No model in this spec! Skipped"
    else
        (cd ${target}/models; go build)
        if [[ $? != 0 ]] ; then
            echo "Build failed for ${spec}"
            if [[ ! -z ${continueOnError} ]] ; then
                failures="${failures} build:${spec}"
                continue
            else
                exit 1
            fi
        fi
        echo "${spec}: Build OK"
        if [[ -n ${clean} ]] ; then
             rm -rf ${target}
        fi
    fi
done
done
if [[ ! -z ${failures} ]] ; then
    echo ${failures}|tr ' ' '\n'
else
    echo "No failures"
fi
exit
# Non reg codegen
# NOTE(fredbi):
# - azure: invalid spec
# - bitbucket: model does not compile
# - issue72: invalid spec
# - todolist.discriminator: ok now
testcases=`cd ../../codegen;ls -1|grep -vE 'azure|bitbucket|existing-model|issue72|todolist.simple.yml'`
#testcases=${testcases}" fixture-1062.json fixture-984.yaml"
#testcases=`cd ../../codegen;ls -1|grep  'todolist.enums.yml'`
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
    echo "Server generation for ${spec}"
    swagger generate server --skip-validation --spec ${spec} --target ${target} --name=${serverName} 1>${testcase%.*}.log 2>&1
    #--log-output=${testcase%.*}.log
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
# TODO(fredbi): enable non reg again
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
