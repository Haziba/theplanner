#!/bin/bash

set -e

aws lambda update-function-code \
    --region us-east-1 \
    --function-name theplanner-food-meals  \
    --zip-file fileb://endpoints/food/meals/deployment.zip

aws lambda update-function-code \
    --region us-east-1 \
    --function-name theplanner-food-ingredients  \
    --zip-file fileb://endpoints/food/ingredients/deployment.zip

aws lambda update-function-code \
    --region us-east-1 \
    --function-name theplanner-food-planner \
    --zip-file fileb://endpoints/food/planner/deployment.zip

aws lambda update-function-code \
    --region us-east-1 \
    --function-name theplanner-money-logs  \
    --zip-file fileb://endpoints/money/logs/deployment.zip
