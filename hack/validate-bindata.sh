#!/bin/bash
# validates bindata templates are generated correctly.
# bindata is generated using github.com/kevinburke/go-bindata/go-bindata@v3.22.0

help_exit(){
cat <<EOF

Usage: validate-bindata.sh
Notes:
Regenerate bindata using github.com/kevinburke/go-bindata/go-bindata@v3.22.0
Fail when bindata.go is modified after regenreate.
Original generate command is in generator/shared.go
EOF
exit 1;
}

check_bindata_installed(){
    which go-bindata || { echo "go-bindata is not installed."; help_exit; } 
}

check_bindata_version(){
    echo $(go-bindata -v) | grep "3.22.0" || { echo "go-bindata version is not 3.22.0"; help_exit; } 
}

regenerate_templates(){
    go generate ./generator/...
}

check_no_file_changed()
{
    git diff --exit-code --name-only -- ./generator  || { echo "bindata.go not up to date. Please commit bindata.go in current state."; help_exit; }
}

check_bindata_installed

check_bindata_version

regenerate_templates

check_no_file_changed

echo "Success: bindata is valid."