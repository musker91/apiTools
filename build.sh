#!/usr/bin/env bash

# create dist directory
if [[ -d "./dist" || -f "./dist.zip" ]]; then
    rm -rf ./dist ./dist.zip   
fi

mkdir ./dist

# mv data file and directory
mv apiTools     ./dist/
cp -rf config   ./dist/
cp -rf data     ./dist/
cp -rf static   ./dist/
cp -rf views    ./dist/

# mkdir logs directory
mkdir ./dist/logs/

# compression program
zip -r dist.zip ./dist

# remove dist directory
rm -rf ./dist
