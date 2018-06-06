#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# Acknowledgements
testcases="${testcases} fixture-1232.yaml"   
#testcases="${testcases} ../1198/fixture-1198.yaml"      # passed ok
# A small utility to build fixture servers
# Fixtures with models only
#testcases="${testcases} fixture-allOf.yaml"  # this one still reveals bugs
# non reg
#testcases="${testcases} ../../codegen/billforward.discriminators.yml"
#testcases="${testcases} ../1277/cloudbreak.json"

#testcases="${testcases} fixture-1479-part.yaml"        # passed ok
#testcases="${testcases} fixture-simple-allOf.yaml"     # passed ok
#testcases="${testcases} fixture-complex-allOf.yaml"    # passed ok
#testcases="${testcases} fixture-is-nullable.yaml"      # passed ok
#testcases="${testcases} fixture-itching.yaml"          # passed ok
#testcases="${testcases} fixture-additionalProps.yaml"  # passed ok
#testcases="${testcases} fixture-tuple.yaml"            # passed ok
#testcases="fixture-1047.yaml fixture-1093.yaml fixture-1066.yaml fixture-arrays.yaml fixture-nested.yaml fixture-1203.yaml fixture-1328.yaml fixture-x-const.yaml fixture-ci.yaml fixture-types.yaml quay-discovery.json fixture-body.yaml fixture-951.json fixture-1062.json fixture-984.yaml"
#testcases="${testcase} todolist.schemavalidation.yml todolist.enums.yml fixture-tuple.yaml fixture-edge.yaml fixture-1047.yaml fixture-1093.yaml fixture-1066.yaml fixture-arrays.yaml fixture-nested.yaml fixture-1203.yaml fixture-1328.yaml fixture-x-const.yaml fixture-ext.yaml fixture-ci.yaml fixture-types.yaml fixture-body.yaml"
#testcases="${testcases} fixture-simple-tuple.yaml"
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
