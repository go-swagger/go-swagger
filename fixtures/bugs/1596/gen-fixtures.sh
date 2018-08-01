#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
# Fixtures with models only
testcases="${testcases} fixture-1596.yaml"
# For this, skip --with-expand (contains polymorphic models):
#testcases="${testcases} ../../codegen/todolist.models.yml"
for opts in "--with-flatten=minimal" "--with-flatten=full" "--with-expand" ;do
	for testcase in ${testcases} ; do
	    target=./gen-`basename ${testcase%.*}`
        case ${opts} in
            "--with-flatten=minimal") 
                target=${target}"-minimal"
                ;;
            "--with-flatten=full") 
                target=${target}"-full"
                ;;
            "--with-expand") 
                target=${target}"-expand"
                ;;
        esac
	    spec=./${testcase}
	    serverName="codegensrv"
	    rm -rf ${target}
	    mkdir ${target}
	    echo "Model generation for ${spec} with ${opts}"
	    swagger generate server --skip-validation ${opts} --spec ${spec} --target ${target} --name=${serverName} --output=${testcase%.*}.log
	    # 1>x.log 2>&1
	    #
	    if [[ $? != 0 ]] ; then
	        echo "Generation failed for ${spec}"
	        exit 1
	    fi
	    echo "${spec}: Generation OK"
	    (cd ${target}/cmd/${serverName}-server; go build)
	    if [[ $? != 0 ]] ; then
	        echo "Build failed for ${spec}"
	        exit 1
	    fi
	    echo "${spec}: Build OK"
	    if [[ -n ${clean} ]] ; then 
	        rm -rf ${target}
	    fi
	done
done
exit
