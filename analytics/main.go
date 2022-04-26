package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Request struct {
	Alias string `json:"alias"`
}

type MondoObj struct {
	Id     string `bson:"_id"`
	Clicks int    `bson:"clicks"`
}

func Increment(collection *mongo.Collection, ctx context.Context, alias string) primitive.M {
	var data bson.M
	if err := collection.FindOne(ctx, bson.M{"_id": alias}).Decode(&data); err != nil {
		collection.InsertOne(ctx, bson.D{
			{Key: "_id", Value: alias},
			{Key: "clicks", Value: 1},
		})
	} else {
		var s MondoObj
		bsonBytes, _ := bson.Marshal(data)
		bson.Unmarshal(bsonBytes, &s)
		collection.UpdateOne(ctx, bson.M{"_id": alias},
			bson.D{
				{"$set", bson.D{{"clicks", s.Clicks + 1}}},
			},
		)
	}
	return data
}

func Read(collection *mongo.Collection, ctx context.Context, alias string) {
	var data bson.M
	if err := collection.FindOne(ctx, bson.M{"_id": alias}).Decode(&data); err != nil {
		fmt.Println("cannot read")
	} else {
		fmt.Println(data)
	}
}

func main() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://shubham:shubhamgupta@cluster0.ljvzd.mongodb.net/myFirstDatabase?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	analyticsDatabase := client.Database("analytics-database")
	analyticsCollection := analyticsDatabase.Collection("analytics-collection")
	Increment(analyticsCollection, ctx, "google")
	Read(analyticsCollection, ctx, "google")
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://shubham:shubhamgupta@cluster0.ljvzd.mongodb.net/myFirstDatabase?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	analyticsDatabase := client.Database("analytics-database")
	analyticsCollection := analyticsDatabase.Collection("analytics-collection")
	resp := events.APIGatewayProxyResponse{
		Headers: make(map[string]string),
	}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	// Parse request body
	rb := Request{}
	if err := json.Unmarshal([]byte(request.Body), &rb); err != nil {
		return resp, err
	}
	Increment(analyticsCollection, ctx, rb.Alias)
	resp.StatusCode = http.StatusOK
	return resp, nil
}
