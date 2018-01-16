#! /bin/bash 
# A small utility to build fixture servers
# for non regression testing of codegen

shopt -s extglob
function okcr() {
tput setaf 2;printf "%s\n" "$*"; tput sgr0
}
function ok() {
tput setaf 2;printf "%s" "$*"; tput sgr0
}
function errcr() {
tput setaf 1 bold;printf "%s\n" "$*"; tput sgr0
}
function err() {
tput setaf 1 bold;printf "%s" "$*"; tput sgr0
}
function warncr() {
tput setaf 3;printf "%s\n" "$*"; tput sgr0
}
function warn() {
tput setaf 3;printf "%s" "$*"; tput sgr0
}
function infocr() {
tput setaf 4;printf "%s\n" "$*"; tput sgr0
}

function info() {
tput setaf 4;printf "%s" "$*"; tput sgr0
}

if [ ! -f `which swagger` ]; then
  echo "can't find swagger in the PATH"
  exit 1
fi

# The following ones fail generation: todo existing_model requires pregeneration
# NOTE(fredbi): 
# - issue72: invalid spec
# - todolist.discriminators.yml: invalid mode generation (issue#1376)
# - bitbucket.json: model does not compile (anonymous allOf1)
# - azure-text-analyis.json: invalid specification with duplicate property in allOf construct (provided fixed version for testing)
# - todolist.simple.yml: invalid default values put on purpose for UT (provided fixed version for testing)
known_failed="@(\
azure-text-analyis.json|\
bitbucket.json|\
existing-model.yml|\
issue72.json|\
todolist.discriminators.yml|\
todolist.simple.yml\
)"
# The following ones should fail validation, but produce correct generated code
known_skip_validation="@(\
todolist.enums.yml|\
todolist.models.yml|\
todolist.schemavalidation.yml\
)"
# All fixtures in ./fixtures/codegen + some others
cd ${0%%/*}
specdir="../fixtures/codegen ../fixtures/bugs/909"
gendir=./tmp-gen
rm -rf ${gendir}

check_list=`for d in ${specdir}; do ls $d/*.yml;ls $d/*.json;ls $d/*.yaml;done 2>/dev/null`

for spec in ${check_list}; do 
    testcase=${spec##*/}

    case ${testcase} in 
    ${known_failed})
        warncr "[`date +%T`]${spec}: not tested against full build because of known issues."
        run="false"
        ;;
    ${known_skip_validation})
        info "[`date +%T`]${spec}: assumed invalid but tested against full build..."
        run="true"
        opts="--skip-validation"
        ;;
    *)
        printf "[`date +%T`]%s: %s..." ${spec} "assumed valid and tested against build"
        run="true"
        opts=""
        ;;
    esac

    if [[ ${run} == "true" ]]; then
        target=${gendir}/gen-${testcase%.*}
        server_name="nrcodegen"
        errlog=${gendir}/stderr.log

        rm -rf ${target}
        mkdir -p ${target}
        rm -f ${errlog}

        # Gen server
        swagger generate server --spec ${spec} --target ${target} --name=${server_name} --quiet ${opts} 2>${errlog}
        if [[ $? != 0 ]] ; then
            errcr "Generation Failed"
            if [[ -f ${errlog} ]] ; then errcr `cat ${errlog}` ; rm ${errlog};fi
            exit 1
        fi
        ok `printf "%s..."  "Generation OK"`
        rm -f ${errlog}
        # Build server
        (cd ${target}/cmd/${server_name}"-server"; go build) 2>${errlog}
        if [[ $? != 0 ]] ; then
            errcr "Build Failed"
            if [[ -f ${errlog} ]] ; then errcr `cat ${errlog}` ; rm ${errlog};fi
            exit 1
        fi
        okcr "Build: OK"
        rm -f ${errlog}
    fi
done
rm -rf ${gendir}
exit 0
