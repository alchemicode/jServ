#!/bin/bash

cp -rv Server TESTING

pushd TESTING
rm Databases/db.json

dart main.dart

popd