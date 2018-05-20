#!/bin/bash

set -e

declare -a endpoints=("food/meals" "food/meals/meal-id" "food/ingredients" "food/ingredients/ingredient-id" "food/planner/build" "food/planner/save" "money/logs" "money-event-sourced/money-change")

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
