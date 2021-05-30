#! /bin/bash

BASEDIR="$GOPATH/src/go-kakeibo"

bash $BASEDIR/.codegen/internal/xo/codegen.sh
bash $BASEDIR/.codegen/cmd/api/openapi/codegen.sh
bash $BASEDIR/.codegen/cmd/api/xo/codegen.sh
bash $BASEDIR/.codegen/test/codegen.sh