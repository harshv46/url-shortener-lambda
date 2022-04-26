package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	LinksTableName = "UrlShortenerLinks"
	Region         = "ap-south-1"
)

type Link struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

// Start DynamoDB session
var sess, sess_err = session.NewSession(&aws.Config{
	Region: aws.String(Region),
})

var svc = dynamodb.New(sess)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get short_url parameter
	shortURL, _ := request.PathParameters["short_url"]
	// If session not started correctly, try again
	if sess_err != nil {
		sess, sess_err := session.NewSession(&aws.Config{
			Region: aws.String(Region),
		})
		if sess_err != nil {
			return events.APIGatewayProxyResponse{}, sess_err
		}
		svc = dynamodb.New(sess)
	}
	// Read link
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(LinksTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String(shortURL),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Unmarshal link item
	link := Link{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &link); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Redirect to long URL
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"location": link.LongURL,
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
