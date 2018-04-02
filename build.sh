#!/bin/bash

set -e

declare -a endpoints=("meals")

for i in "${endpoints[@]}"
do
	cd endpoints/$i
	GOOS=linux go build -o main

	if [[ $(uname) == "MINGW"* ]]; then
		~/go/bin/build-lambda-zip.exe -o deployment.zip main
	else
		zip deployment.zip main
	fi

	cd -
done
