#!/usr/bin/env bash

#  create dist directory
if [[ -d "./dist" ]]; then
    echo "dist directory is exists, please first delete!!!"
    exit
else
    mkdir dist
fi

# go build linux amd64
GOOS=linux GOARCH=amd64 go build -o apiTools main.go
if [[ $? != 0 ]];then
    echo "go build fail!!!"
    exit
fi

# mv data file and directory
mv apiTools ./dist/
cp -rf config ./dist/
cp -rf data   ./dist/
cp -rf static  ./dist/
cp -rf views  ./dist/

# mkdir logs directory
mkdir ./dist/logs/

# compression program
zip -r dist.zip ./dist

# remove dist directory
rm -rf ./dist