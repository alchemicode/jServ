#!/bin/bash

rm -rfv TESTING

cp -rv Server TESTING

pushd TESTING
rm Databases/db.json
echo "[]" > Databases/SparkPlug.json
sed -i 's/admin/SparkPlug/g' data.jserv

dart main.dart

popd