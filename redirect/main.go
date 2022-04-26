package main

import (
	"net/http"
	"time"
	// "fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-redis/redis/v7"
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
var redisClient = redis.NewClient(&redis.Options{
		Addr: "lambdacache.ncfh49.0001.aps1.cache.amazonaws.com:6379",
		DB: 0,
	})

func CacheSet(key string, post string) {

	redisClient.Set(key, post, time.Duration(3600)*time.Second)
}

func CacheGet(key string) (string, error) {

	val, err := redisClient.Get(key).Result()
	if err != nil {
		return "",nil
	}

	return val,nil
}

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
	link := Link{}
	// Read link
	// in := "out"
	longURL, _ := CacheGet(shortURL)
	if longURL == "" {
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
		if err := dynamodbattribute.UnmarshalMap(result.Item, &link); err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		// in = "in"
		CacheSet(shortURL, link.LongURL)
	} else {
		link.ShortURL = shortURL
		link.LongURL = longURL
	}
	// Redirect to long URL
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"location": link.LongURL,
		},
	}, nil
}


// func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)  {
// 	fmt.Println("Go Redis Tutorial")

// 	pong, err := redisClient.Ping().Result()
// 	fmt.Println(pong, err)

// 	CacheSet("test", "val")

// 	val, err := CacheGet("test1")
// 	in := "out"
// 	if val == "" {
// 		in = "in"
// 	fmt.Println(err)
// 	}

// 	fmt.Println(val)

// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusPermanentRedirect,
// 		Headers: map[string]string{
// 			"location": in,
// 		},
// 		Body: "",
// 	}, nil

// }

func main() {
	lambda.Start(Handler)
}
