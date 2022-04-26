GOOS=linux GOARCH=amd64 go build -o redirect main.go
zip deployment.zip analytics

aws lambda update-function-code --function-name AnalyticsFunction --region ap-south-1 --zip-file fileb://./deployment.zip
