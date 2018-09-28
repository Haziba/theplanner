#!/bin/bash

set -e

aws cloudformation deploy --region eu-west-1 --template-file ./template.yaml --stack-name theplanner

#aws lambda update-function-code \
#    --region us-east-1 \
#    --function-name theplanner-food-meals  \
#    --zip-file fileb://endpoints/food/meals/deployment.zip


#aws lambda update-function-code \
#    --region us-east-1 \
#    --function-name theplanner-food-ingredients  \
#    --zip-file fileb://endpoints/food/ingredients/deployment.zip

#aws lambda update-function-code \
#    --region us-east-1 \
#    --function-name theplanner-food-planner-build \
#    --zip-file fileb://endpoints/food/planner/build/deployment.zip

#aws lambda update-function-code \
#    --region us-east-1 \
#    --function-name theplanner-food-planner-save \
#    --zip-file fileb://endpoints/food/planner/save/deployment.zip

#aws lambda update-function-code \
#    --region us-east-1 \
#    --function-name theplanner-money-logs  \
#    --zip-file fileb://endpoints/money/logs/deployment.zip

#aws lambda update-function-code \
#    --region us-east-1 \
#    --function-name theplanner-money-logs-2  \
#    --zip-file fileb://endpoints/money-event-sourced/change-money/deployment.zip
