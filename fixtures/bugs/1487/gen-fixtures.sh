#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
continueOnError=1
# A small utility to build fixture servers
# Fixtures with models only


#testcases="${testcases} ../../codegen/issue72.json"
#testcases="${testcases} ../../canary/bitbucket.org/swagger.json"
#testcases="${testcases} ../../codegen/azure-text-analyis.json"
#testcases="${testcases} ../../codegen/todolist.simple.yml"
#testcases="${testcases} swagger-gsma.yaml"
#testcases="${testcases} ../844/swagger.json"
#testcases="${testcases} fixture-844-variations.yaml"
#testcases="${testcases} fixture-anonymous-enum.yaml"
#testcases="${testcases} fixture-nested-maps.yaml"
#testcases="${testcases} fixture-errors.yaml"
#testcases="${testcases} fixture-simple-allOf.yaml"
#testcases="${testcases} fixture-complex-allOf.yaml"
#testcases="${testcases} fixture-is-nullable.yaml"
#testcases="${testcases} fixture-itching.yaml"
#testcases="${testcases} fixture-additionalProps.yaml"
#testcases="${testcases} fixture-tuple.yaml"
#testcases="${testcases} fixture-allOf.yaml"
testcases="${testcases} ../1479/fixture-1479-part.yaml"
#testcases="${testcases} ../1198/fixture-1198.yaml"
#testcases="${testcases} ../1042/fixture-1042.yaml"
#testcases="${testcases} ../1042/fixture-1042-2.yaml"
#testcases="${testcases} ../979/fixture-979.yaml"
#testcases="${testcases} ../842/fixture-842.yaml"
#testcases="${testcases} ../607/fixture-607.yaml"
#testcases="${testcases} ../1336/fixture-1336.yaml"
#testcases="${testcases} ../1277/cloudbreak.json"
#testcases="${testcases} ../../codegen/todolist.schemavalidation.yml"
#testcases="${testcases} ../../codegen/todolist.discriminators.yml"
#testcases="${testcases} ../../codegen/billforward.discriminators.yml"
#opts="--skip-flatten"
#export SWAGGER_DEBUG=1
for testcase in ${testcases} ; do
    spec=${testcase}
    testcase=`basename ${testcase}`
    if [[ -z ${opts} ]]; then
        target=./gen-${testcase%.*}-flatten
    else
        target=./gen-${testcase%.*}-expand
    fi
    serverName="codegensrv"
    rm -rf ${target}
    mkdir ${target}
    echo "Model generation for ${spec}"
    swagger generate model --skip-validation ${opts} --spec=${spec} --target=${target} --output=${testcase%.*}.log
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
