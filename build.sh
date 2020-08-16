#!/usr/bin/env bash

# create dist directory
if [[ -d "./dist" || -f "./dist.tar.gz" ]]; then
    rm -rf ./dist ./dist.tar.gz
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
tar zcf dist.tar.gz ./dist

# remove dist directory
rm -rf ./dist
