#aws cloudformation delete-stack --stack-name theplanner-api --region us-east-1
#aws cloudformation stack-delete-complete --stack-name theplanner-api --region us-east-1
aws cloudformation deploy --stack-name theplanner-api --template-file template.yaml --parameter-overrides ParameterKey=thing,ParameterValue=cool --region us-east-1