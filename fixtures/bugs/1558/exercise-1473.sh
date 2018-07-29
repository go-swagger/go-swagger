#! /bin/bash
fixturePath="./gen-fixture-1558-flatten/cmd/nrcodegen-server"
cmd="${fixturePath}/nrcodegen-server"
# Generates a fresh server
./gen-fixtures.sh
# Apply wrong options like in #1473
patch gen-fixture-1558-flatten/restapi/configure_nrcodegen.go patch-1473-1.patch
#
(cd ${fixturePath} ; go build)
if [[ ! $? -eq 0 ]] ; then
    echo "Error building server"
    exit 1
fi
captured="./srv.log"
# This should panic and not hang
$cmd --help 1>${captured} 2>&1
grep -qi "panic: provided data is not a pointer to struct" ${captured}
if [[ $? -eq 0 ]] ; then
    echo "Server startup panicked as expected. OK"
else
    echo "Server startup did not panick as expected. Wrong. Here is the output:"
    echo ">>>"
    cat ${captured}
    echo ">>>"
fi

# Apply patch to correct this
echo "Now fixing the issue..."
patch gen-fixture-1558-flatten/restapi/configure_nrcodegen.go patch-1473-2.patch
(cd ${fixturePath} ; go build)
# This should be ok now
echo "Try again..."
$cmd --help 1>${captured} 2>&1
grep -qi "example1" ${captured}
if [[ $? -eq 0 ]] ; then
    echo "Server startup shows help with additional option. OK"
else
    echo "Server startup did not work as expected. Wrong. Here is the output:"
    echo ">>>"
    cat ${captured}
    echo ">>>"
fi
