#!/bin/bash

# set -e

aws --region us-east-1 dynamodb create-table \
    --table-name tp-meals \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

aws --region us-east-1 dynamodb create-table \
    --table-name tp-ingredients \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

aws --region us-east-1 dynamodb create-table \
    --table-name tp-money-logs \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

aws lambda create-function \
    --region us-east-1 \
    --function-name theplanner-food-meals \
    --memory 128 \
    --role arn:aws:iam::874444545532:role/service-role/theplanner \
    --runtime go1.x \
    --zip-file fileb://endpoints/food/meals/deployment.zip \
    --handler main \
    --environment Variables="{PROD=true}"

aws lambda create-function \
    --region us-east-1 \
    --function-name theplanner-food-ingredients \
    --memory 128 \
    --role arn:aws:iam::874444545532:role/service-role/theplanner \
    --runtime go1.x \
    --zip-file fileb://endpoints/food/ingredients/deployment.zip \
    --handler main \
    --environment Variables="{PROD=true}"

aws lambda create-function \
    --region us-east-1 \
    --function-name theplanner-food-planner \
    --memory 128 \
    --role arn:aws:iam::874444545532:role/service-role/theplanner \
    --runtime go1.x \
    --zip-file fileb://endpoints/food/planner/deployment.zip \
    --handler main \
    --environment Variables="{PROD=true}"

aws lambda create-function \
    --region us-east-1 \
    --function-name theplanner-money-logs \
    --memory 128 \
    --role arn:aws:iam::874444545532:role/service-role/theplanner \
    --runtime go1.x \
    --zip-file fileb://endpoints/money/logs/deployment.zip \
    --handler main \
    --environment Variables="{PROD=true}"
