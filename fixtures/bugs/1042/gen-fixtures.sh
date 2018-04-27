#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# Acknowledgements
testcases="${testcases} fixture-1042.yaml"   
testcases="${testcases} fixture-1042-2.yaml"   
for testcase in ${testcases} ; do
    spec=${testcase}
    testcase=`basename ${testcase}`
    target=./gen-${testcase%.*}
    serverName="codegensrv"
    rm -rf ${target}
    mkdir ${target}
    echo "Model generation for ${spec}"
    swagger generate model --skip-validation --spec ${spec} --target ${target} --output=${testcase%.*}.log
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
        (cd ${target}/models; go build)
        if [[ $? != 0 ]] ; then
            echo "Build failed for ${spec}"
            exit 1
        fi
        echo "${spec}: Build OK"
        if [[ -n ${clean} ]] ; then 
             rm -rf ${target}
        fi
    fi
done
exit
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
    #--output=${testcase%.*}.log
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
