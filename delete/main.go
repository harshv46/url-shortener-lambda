package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	LinksTableName = "UrlShortenerLinks"
	Region         = "ap-south-1"
)

type Link struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get short_url parameter
	var shortURL, _ = request.PathParameters["short_url"]
	// Start DynamoDB session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	svc := dynamodb.New(sess)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String(shortURL),
			},
		},
		TableName: aws.String("UrlShortenerLinks"),
	}
	_, err = svc.DeleteItem(input)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Redirect to long URL
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
		Headers: map[string]string{
			"deleted": shortURL,
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
