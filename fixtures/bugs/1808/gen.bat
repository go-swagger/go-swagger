del /S /F /Q tmp && mkdir tmp
%GOPATH%\bin\swagger.exe generate server -f fixture-1808.yaml /A vg-api /m restapimodels --target=tmp
