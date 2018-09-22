fixture=fixture-1084.yaml
verifier=unmarshal_test.go
target=gen-fixture-1084
if [[ ! -f $verifier ]] ; then
    mv ${verifier%.go}.gol ${verifier%.go}.go
fi
if [[ ! -d $target ]] ; then
    mkdir $target
fi
swagger generate model --spec=$fixture --target=$target --quiet
go test -vet off -v
ret=$?
if [[ $1 == "--clean" ]] ; then
    rm -rf $target
    mv ${verifier%.go}.go ${verifier%.go}.gol
fi
exit $ret
