#! /bin/bash

# xo オプション
# -1 : toggle query's generated Go func to return only one result
# -M : toggle trimming of query whitespace in generated Go code
# -N : enable query mode
# -T : query's generated Go type
# -o : output path or file name
# -p : package name used in generated Go code
# --template-path : user supplied template path

cd `dirname $0`

rm -rf out
mkdir -p out/entity out/repository

DSN="mysql://root:password@mysql-server/kakeibo?parseTime=true"
xo $DSN -o out/entity --template-path templates/entity
xo $DSN -o out/repository --template-path templates/repository
