#! /bin/bash

cd `dirname $0`

BASEDIR=$GOPATH/src/go-kakeibo
find $BASEDIR -name mock_*.go -exec rm {} \;
mockery --all --inpackage --dir $BASEDIR 
