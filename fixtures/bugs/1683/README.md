# A few generations tested

Below find the 5424 combinations tested of `--model-package`, `--api-package`, `--target`, `--client-package` and `--name`.
Tests involve: using sub-directories, putting dashes and underscores in the package name.

``` 
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc --server-package=sabc/sub1/ssub2 --api-package=aabc/asubdir --target=codegen2/target-withDash264 --name=nrcodegen
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc --server-package=sabc/ssubdir --api-package=aabc-dashed --target=codegen2/target-withDash209 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc/sub1/ssub2 --api-package=aabc/asubdir --target=codegen2/target-withDash581 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=sabc-dashed/ssubdir --api-package=aabc-dashed --target=codegen2/target-withDash871 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc-dashed --api-package=aabc-dashed --target=codegen2/target-withDash540 --name=nrcodegen-test
swagger generate client --skip-validation --spec=fixture-1683.yaml --model-package=models --client-package=cabc/csubdir --target=codegen2/target-withDash151 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc --server-package=sabc/ssubdir --api-package=operations --target=codegen2/target-withDash206 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed/msubdir --server-package=sabc/ssubdir --api-package=aabc-dashed --target=codegen2/target-withDash667 --name=nrcodegen
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed/msubdir --server-package=restapi --api-package=aabc --target=codegen2/target-withDash623 --name=nrcodegen
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=sabc --api-package=aabc/asubdir --target=codegen2/target-withDash814 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/msubdir --server-package=sabc-test --api-package=aabc-dashed/asubdir --target=codegen2/target-withDash456 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc-dashed --api-package=aabc --target=codegen2/target-withDash538 --name=nrcodegen-test
swagger generate client --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed/msubdir --client-package=cabc-dashed --target=codegen2/target-withDash772 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc-dashed --api-package=aabc-test --target=codegen2/target-withDash536 --name=nrcodegen
...
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir507 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=restapi --api-package=aabc-dashed/asubdir --target=codegen3/target/subDir795 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/msubdir --server-package=sabc/sub1/ssub2 --api-package=aabc-test --target=codegen3/target/subDir430 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=sabc-dashed --api-package=aabc-test --target=codegen3/target/subDir853 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=models --server-package=sabc/sub1/ssub2 --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir112 --name=nrcodegen
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=models --server-package=sabc-dashed/ssubdir --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir91 --name=nrcodegen
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc-test --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir612 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/msubdir --server-package=sabc-dashed/ssubdir --api-package=aabc/asubdir --target=codegen3/target/subDir405 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=sabc --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir810 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir500 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=sabc/sub1/ssub2 --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir901 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-test --server-package=restapi --api-package=aabc --target=codegen3/target/subDir947 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc --api-package=aabc-test --target=codegen3/target/subDir501 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-test --server-package=sabc/ssubdir --api-package=aabc-dashed/asubdir --target=codegen3/target/subDir992 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=models --server-package=sabc-test --api-package=operations --target=codegen3/target/subDir135 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/msubdir --server-package=restapi --api-package=aabc-dashed/asubdir --target=codegen3/target/subDir330 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed/msubdir --server-package=restapi --api-package=aabc-dashed/asubdir --target=codegen3/target/subDir640 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc/ssubdir --api-package=aabc-dashed --target=codegen3/target/subDir512 --name=nrcodegen
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed/msubdir --server-package=sabc-dashed/ssubdir --api-package=aabc/asubdir --target=codegen3/target/subDir722 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed/msubdir --server-package=sabc/sub1/ssub2 --api-package=aabc-test --target=codegen3/target/subDir740 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-test --server-package=sabc/ssubdir --api-package=aabc-test --target=codegen3/target/subDir987 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc-test --api-package=operations --target=codegen3/target/subDir600 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=sabc-dashed --api-package=aabc-dashed/asubdir --target=codegen3/target/subDir844 --name=nrcodegen
swagger generate client --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed/msubdir --client-package=cabc --target=codegen3/target/subDir770 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc-dashed/ssubdir --api-package=aabc --target=codegen3/target/subDir559 --name=nrcodegen-test
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=sabc-test --api-package=operations --target=codegen3/target/subDir903 --name=nrcodegen
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc-dashed --server-package=sabc-test --api-package=aabc-test --target=codegen3/target/subDir613 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=models --server-package=sabc --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir42 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/sub1/msub2 --server-package=sabc/ssubdir --api-package=aabc/sub1/asub2 --target=codegen3/target/subDir824 --name=nrcodegen
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc/msubdir --server-package=sabc-test --api-package=aabc-test --target=codegen3/target/subDir458 --name=nrcodegen_underscored
swagger generate server --skip-validation --spec=fixture-1683.yaml --model-package=mabc --server-package=sabc-dashed/ssubdir --api-package=aabc --target=codegen3/target/subDir256 --name=nrcodegen_underscored
```
