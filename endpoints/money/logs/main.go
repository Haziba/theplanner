package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	helpers "github.com/haziba/theplanner/helpers"
	models "github.com/haziba/theplanner/models/money"
	money "github.com/haziba/theplanner/services/money/log"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	moneyLogService, err := helpers.CreateMoneyLogService()
	if err != nil {
		log.Printf("error creating money log service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	if request.HTTPMethod == "POST" {
		return post(request, moneyLogService)
	}

	return get(moneyLogService)
}

func post(request events.APIGatewayProxyRequest, moneyLogService money.MoneyLogService) (events.APIGatewayProxyResponse, error) {
	var logs []models.MoneyLog

	err := json.Unmarshal([]byte(request.Body), &logs)

	if err != nil {
		log.Printf("error unmarshalling ingredient: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	for _, moneyLog := range logs {
		moneyLog, err = moneyLogService.CreateMoneyLog(moneyLog)
		if err != nil {
			log.Printf("error creating log: %v\n", err)
			return helpers.CreateBadRequestResponse()
		}
	}

	data, err := json.Marshal(logs)
	if err != nil {
		log.Printf("error marshalling log: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(data),
		StatusCode: 200,
	}, nil
}

func get(moneyLogService money.MoneyLogService) (events.APIGatewayProxyResponse, error) {
	logs, err := moneyLogService.GetAllMoneyLogs()

	if err != nil {
		log.Printf("error getting logs %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	m := moneyLogResponse{
		Logs: logs,
	}

	data, err := json.Marshal(m)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(data),
		StatusCode: 200,
	}, nil
}

type moneyLogResponse struct {
	Logs []models.MoneyLog `json:"logs"`
}

func createMoneyLogResponse(logs []models.MoneyLog) (events.APIGatewayProxyResponse, error) {
	resp := moneyLogResponse{
		Logs: logs,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(data),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
