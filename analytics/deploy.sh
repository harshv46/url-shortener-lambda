GOOS=linux GOARCH=amd64 go build -o analytics main.go
zip deployment.zip analytics
aws lambda create-function --function-name AnalyticsFunction --region ap-south-1 --zip-file fileb://./deployment.zip --runtime go1.x --tracing-config Mode=Active --role $ROLE --handler analytics
