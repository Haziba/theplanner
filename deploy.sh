#!/bin/bash

set -e

aws lambda update-function-code \
    --function-name theplanner-meals  \
    --zip-file fileb://endpoints/meals/deployment.zip
