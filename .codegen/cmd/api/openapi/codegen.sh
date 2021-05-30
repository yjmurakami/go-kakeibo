#! /bin/bash

# openapi-generator オプション
#
# -i <spec file>, --input-spec <spec file>
#     location of the OpenAPI spec, as URL or file (required if not loaded
#     via config using -c)
# 
# -g <generator name>, --generator-name <generator name>
#     generator to use (see list command for list)
# 
# -o <output directory>, --output <output directory>
#     where to write the generated files (current dir by default)
# 
# -t <template directory>, --template-dir <template directory>
#     folder containing the template files
# 
# --ignore-file-override <ignore file override location>
#     Specifies an override location for the .openapi-generator-ignore
#     file. Most useful on initial generation.
# 
# --type-mappings <type mappings>
#     sets mappings between OpenAPI spec types and generated code types in
#     the format of OpenAPIType=generatedType,OpenAPIType=generatedType.
#     For example: array=List,map=Map,string=String. You can also have
#     multiple occurrences of this option.

cd `dirname $0`

rm -rf out

java -jar $GOPATH/bin/openapi-generator-cli.jar generate \
    -i $GOPATH/src/go-kakeibo/openapi.yml \
    -g go-server \
    -o out \
    -t templates \
    --ignore-file-override templates/.openapi-generator-ignore \
    --type-mappings integer=int

rename .go .oa.go out/go/*.go

# go fmtを実行するためにpackageの不整合を解消する
rm out/go/model_inline_*
mkdir out/go/model
mv out/go/model_* out/go/model/
rename model_ "" out/go/model/*.go

mv out/go/api.oa.go out/go/handler.oa.go
mv out/go/routers.oa.go out/go/router.oa.go

sed -i -e 's/ApiHandler/Handler/g' out/go/handler.oa.go
sed -i -e 's/ApiHandler/Handler/g' out/go/router.oa.go

go fmt ./...