#! /bin/bash
fixturePath="./gen-fixture-1558-flatten/cmd/nrcodegen-server"
cmd="${fixturePath}/nrcodegen-server"
captured="./srv.log"

# build fresh server
./gen-fixtures.sh

function runWithOpts() {
    opts=$1
    expectedRet=${2:-0}
    echo "Running with opts: ${opts}"
    echo "Expected return: ${expectedRet}"

    echo "${cmd} ${opts} 1>${captured} 2>&1 &"
    ${cmd} ${opts} 1>${captured} 2>&1 &
    pid=$!
    sleep 1
    ps -efo pid|grep -q ${pid}
    isInactive=$? 
    if [[ ${isInactive} -eq 0 ]] ; then
        ./killme.sh ${pid}
    fi
    wait $pid
    ret=$?
    if [[ ${ret} -eq 0 ]] ; then
        if [[ ! ${expectedRet} -eq 0 ]] ; then
            echo "Expected startup error, got a success instead"
            wrong=1
        fi
        kill ${pid}
    else
        if [[ ${expectedRet} -eq 0 ]] ; then
            echo "Expected correct startup, got an error instead"
            wrong=1
        else
            echo "Expected startup error. OK"
        fi
    fi
    echo "Terminated with ret=${ret}"
    if [[ ${expectedRet} -eq 0 ]] ; then
        # verifies graceful shutdown
        grep -qi "Shutting down" ${captured}
        cond1=$?
        grep -qi "Stopped serving" ${captured}
        cond2=$?
        if [[ ${cond1} -ne 0 || ${cond2} -ne 0 ]] ; then
            echo "expected graceful shutdown, got:"
            wrong=1
        else
            echo "Terminated gracefully. OK"
        fi
    fi
    if [[ ${wrong} -eq 1 ]] ; then 
        echo "Unexpected startup behavior"
    fi
    echo "Here is the captured output:"
    cat ${captured}
    rm -f ${captured}
}

opts1="\
    --scheme=https \
    --tls-host=0.0.0.0 \
    --tls-port=12345 \
    --tls-certificate=mycert1.crt \
    --tls-ca=myCA.crt \
    --tls-key=mycert1.key"

runWithOpts "${opts1}"

# wrong ca cert
opts2="\
    --scheme=https \
    --tls-host=0.0.0.0 \
    --tls-port=12345 \
    --tls-certificate=mycert1.crt \
    --tls-ca=noWheremyCA.crt \
    --tls-key=mycert1.key"

runWithOpts "${opts2}" "1"

# encrypted priv key
opts3="\
    --scheme=https \
    --tls-host=0.0.0.0 \
    --tls-port=12345 \
    --tls-certificate=mycert1.crt \
    --tls-key=mycert1.encrypted.key"

runWithOpts "${opts3}" "1"

# wrong cert file
opts4="\
    --scheme=https \
    --tls-host=0.0.0.0 \
    --tls-port=12345 \
    --tls-certificate=noWheremycert1.crt \
    --tls-key=mycert1.key"

runWithOpts "${opts4}" "1"

# corrupted CA cert file
opts5="\
    --scheme=https \
    --tls-host=0.0.0.0 \
    --tls-port=12345 \
    --tls-certificate=mycert1.crt \
    --tls-ca=myCA.corrupted.crt \
    --tls-key=mycert1.key"

runWithOpts "${opts5}" "1"

# corrupted cert file
opts6="\
    --scheme=https \
    --tls-host=0.0.0.0 \
    --tls-port=12345 \
    --tls-certificate=mycert1.corrupted.crt \
    --tls-key=mycert1.key"

runWithOpts "${opts6}" "1"

# missing cert file
opts7="\
    --scheme=https \
    --tls-host=0.0.0.0 \
    --tls-port=12345 \
    --tls-key=mycert1.key"

runWithOpts "${opts7}" "1"

# missing key file
opts8="\
    --scheme=https \
    --tls-host=0.0.0.0 \
    --tls-port=12345 \
    --tls-certificate=mycert1.crt"

runWithOpts "${opts8}" "1"

