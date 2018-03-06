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
# The following ones should fail validation, but produce correct generated code (at least it builds)
known_skip_validation="@(\
todolist.enums.yml|\
todolist.models.yml|\
todolist.schemavalidation.yml|\
swagger-gsma.json\
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
specdir="../fixtures/codegen ../fixtures/bugs/909 ../fixtures/bugs/1437"
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
        rm -rf ${target}
    fi
done
rm -rf ${gendir}
exit 0
