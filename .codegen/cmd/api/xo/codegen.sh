#! /bin/bash

# xo オプション
# -1 : toggle query's generated Go func to return only one result
# -M : toggle trimming of query whitespace in generated Go code
# -N : enable query mode
# -T : query's generated Go type
# -F : query's generated Go func name
# -o : output path or file name
# -p : package name used in generated Go code
# --query-type-comment : comment for query's generated Go type
# --query-func-comment : comment for query's generated Go func
# --template-path : user supplied template path

cd `dirname $0`

rm -rf out
mkdir out

DSN="mysql://root:password@mysql-server/kakeibo?parseTime=true"
FILES=$(ls sql/*.sql)
for FILE in $FILES
  do
    OPTION=$(head -n 2 $FILE | tail -n 1)         # 2行目
    TYPE_NAME=$(head -n 4 $FILE | tail -n 1)      # 4行目
    TYPE_COMMENT=$(head -n 6 $FILE | tail -n 1)   # 6行目
    FUNC_NAME=$(head -n 8 $FILE | tail -n 1)      # 8行目
    FUNC_COMMENT=$(head -n 10 $FILE | tail -n 1)  # 10行目

    # 12行目以降はSQL
    tail -n +12 $FILE | xo $DSN $OPTION -N -T $TYPE_NAME -F $FUNC_NAME --query-type-comment "$TYPE_COMMENT" --query-func-comment "$FUNC_COMMENT" -o out --template-path templates
done
