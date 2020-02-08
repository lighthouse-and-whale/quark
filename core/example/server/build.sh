#!/usr/bin/env bash
# ${PWD##*/}
export GOPATH=$(pwd):$(cd ../vendor-go;pwd)

version='0.0'
small_version=$(($(cat version)+1)); printf ${small_version} > version
echo """package version
const Version = ${small_version}
""" > src/version/version.go

build_name='engine'
dist_name=${build_name}_v${version}.${small_version}

rm -rf bin
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o bin/${build_name}
cp config.json bin
cp run.sh bin

if [ ! -d 'dist' ]; then mkdir dist; fi
cd bin; tar -cf - * | pigz  > ../dist/${dist_name}_linux_amd64.tar.gz
