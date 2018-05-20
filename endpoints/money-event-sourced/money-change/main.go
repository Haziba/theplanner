package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/haziba/theplanner/helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	S3_REGION = "us-east-1"
	S3_BUCKET = "theplanner-streams"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	conf := aws.Config{Region: aws.String(S3_REGION)}
	sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)

	filename := "coolfile.txt"

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Failed to open file", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	log.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(S3_BUCKET),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})
	if err != nil {
		log.Println("error", err)
		os.Exit(1)
	}

	log.Printf("Successfully uploaded %s to %s\n", filename, result.Location)

	data, err := json.Marshal("Sup")
	if err != nil {
		log.Printf("error marshalling log: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(data),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
