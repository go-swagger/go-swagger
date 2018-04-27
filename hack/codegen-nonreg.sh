#! /bin/bash
# A small utility to build fixture servers
# for non regression testing of codegen
shopt -s extglob

function initColors() {
if [[ -z ${MONO} ]]; then
    normal=$(tput sgr0)
    red=$(tput setaf 1)
    green=$(tput setaf 2)
    orange=$(tput setaf 3)
    bold=$(tput bold)
    blue=$(tput setaf 4)
else
    normal=""
    red=""
    green=""
    orange=""
    blue=""
    bold=""
fi
}

function okcr() {
    ok "$*"; printf "\n"
}
function ok() {
    printf "${green}%s${normal}" "$*"
}
function errcr() {
    err "$*" ; printf "\n"
}
function err() {
    printf "${red}${bold}%s${normal}" "$*"
}
function successcr() {
    success "$*" ; printf "\n"
}
function success() {
    printf "${green}${bold}%s${normal}" "$*"
}
function warncr() {
    warn "$*" ; printf "\n"
}
function warn() {
    printf "${orange}%s${normal}" "$*"
}
function infocr() {
    info "$*" ; printf "\n"
}
function info() {
    printf "${blue}%s${normal}" "$*"
}

if [ ! -f `which swagger` ]; then
  echo "can't find swagger in the PATH"
  exit 1
fi

# NOTE(fredbi): 
# The following ones fail generation: 
# - todo existing_model requires pregeneration
# - issue72: model works with --skip-validation. Invalid spec (duplicate operationID)
# - todolist.simple.yml: invalid default values put on purpose for UT (provided fixed version for testing)
#
# The following ones requires some checks to be skipped:
# - azure-text-analyis.json: works with --skip-validation. Invalid specification with duplicate property in allOf construct (provided fixed version for testing)
# - azure-text-analyis.json: should work with -skip-validation
#
# The following ones used to fail and are ok:
# - todolist.discriminators.yml: works
# - bitbucket.json: works nows
# - cloudbreak.json : now works
known_failed="@(\
existing-model.yml|\
issue72.json|\
todolist.simple.yml\
)"
# The following ones should fail validation, but produce correct generated code (at least it builds)
known_skip_validation="@(\
todolist.enums.yml|\
todolist.enums.flattened.json|\
todolist.models.yml|\
todolist.schemavalidation.yml|\
azure-text-analyis.json|\
swagger-gsma.json\
)"

# A list of known client build failures
known_client_failure="@(\
todolist.arrayform.yml|\
todolist.arrayquery.yml|\
todolist.url.basepath.yml|\
todolist.url.simple.yml|\
swagger-codegen-tests.json|\
fixture-1414.json|\
fixture-909-3.yaml|\
fixture-909-4.yaml|\
fixture-909-5.yaml|\
fixture-909-6.yaml|\
gentest2.yaml|\
gentest3.yaml|\
gentest.yaml|\
fixture-1437-4.yaml|\
fixture-1392-2.yaml|\
fixture-1392-3.yaml|\
)"

if [[ "$1" = "--circleci" ]] ; then
    # Coloured output not supported by default on CircleCI.
    # Forcing term to xterm is not enough: tput not available with minimalist env.
    MONO=1
    #export TERM=xterm
    #MONO=""
else
    MONO=""
fi

# A little chrome does not hurt...
initColors

# All fixtures in ./fixtures/codegen + some others
cd ${0%/*}
specdir="../fixtures/codegen ../fixtures/bugs/909 ../fixtures/bugs/1437 ../fixtures/bugs/1314 ../fixtures/bugs/1062/eve-online-esi.json"
specdir=${specdir}" ../fixtures/bugs/1392"
specdir=${specdir}" ../fixtures/bugs/1277"
gendir=./tmp-gen
rm -rf ${gendir}

check_list=`for d in ${specdir}; do ls $d/*.yml;ls $d/*.json;ls $d/*.yaml;done 2>/dev/null`
for spec in ${check_list}; do 
    testcase=${spec##*/}

    case ${testcase} in 
    ${known_failed})
        warncr "[`date +%T`]${spec}: not tested against full build because of known issues."
        run="false"
        opts=""
        buildClient="false"
        ;;
    ${known_skip_validation})
        info "[`date +%T`]${spec}: assumed invalid but tested against full build..."
        run="true"
        opts="--skip-validation"
        buildClient="true"
        ;;
    ${known_client_failure})
        warncr "[`date +%T`]${spec}: will not attempt to build the client because of known issues..."
        run="true"
        opts=""
        buildClient="false"
        ;;
    *)
        printf "[`date +%T`]%s: %s..." ${spec} "assumed valid and tested against build"
        run="true"
        opts=""
        buildClient="true"
        ;;
    esac

    if [[ ${run} == "true" ]]; then
        target=${gendir}/gen-${testcase%.*}
        target_client=${gendir}/gen-${testcase%.*}"-client"
        server_name="nrcodegen"
        client_name="nrcodegen"
        errlog=${gendir}/stderr.log

        rm -rf ${target} ${target_client}
        mkdir -p ${target} ${target_client}
        rm -f ${errlog}

        # Gen server
        swagger generate server --spec ${spec} --target ${target} --name=${server_name} --quiet ${opts} 2>${errlog}
        if [[ $? != 0 ]] ; then
            errcr "Generation Failed"
            if [[ -f ${errlog} ]] ; then errcr `cat ${errlog}` ; rm ${errlog};fi
            exit 1
        fi
        ok `printf " %s..."  "Generation server OK"`
        rm -f ${errlog}
        # Gen client
        swagger generate client --spec ${spec} --target ${target_client} --name=${client_name} --quiet ${opts} 2>${errlog}
        if [[ $? != 0 ]] ; then
            errcr "Generation Failed"
            if [[ -f ${errlog} ]] ; then errcr `cat ${errlog}` ; rm ${errlog};fi
            exit 1
        fi
        ok `printf " %s..."  "Generation client OK"`
        # Build server
        (cd ${target}/cmd/${server_name}"-server"; go build) 2>${errlog}
        if [[ $? != 0 ]] ; then
            errcr "Server build Failed"
            if [[ -f ${errlog} ]] ; then errcr `cat ${errlog}` ; rm ${errlog};fi
            exit 1
        fi
        ok `printf " %s..."  "Server build OK"`
        # Build models if any produced 
        if [[ -d ${target}/models ]] ; then 
            (cd ${target}/models ; go build) 2>${errlog}
            if [[ $? != 0 ]] ; then
                errcr "Model build Failed"
                if [[ -f ${errlog} ]] ; then errcr `cat ${errlog}` ; rm ${errlog};fi
                exit 1
            fi
        fi
        ok `printf " %s..."  "Models build OK"`
        # Build client
        if [[ ${buildClient} == "false" ]] ; then
            warn "(no client built)"
            # continue
        else
            (cd ${target_client}/client ; go build) 2>${errlog}
            if [[ $? != 0 ]] ; then
                errcr "Client build Failed"
                if [[ -f ${errlog} ]] ; then errcr `cat ${errlog}` ; rm ${errlog};fi
                exit 1
            fi
            ok `printf " %s..."  "Client build OK"`
        fi
        successcr "[All builds: OK]"
        rm -f ${errlog}
        rm -rf ${target}
    fi
done
rm -rf ${gendir}
exit 0
