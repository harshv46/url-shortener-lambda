GOOS=linux GOARCH=amd64 go build -o delete main.go
zip deployment.zip delete
aws lambda update-function-code --function-name DeleteFunction --region ap-south-1 --zip-file fileb://./deployment.zip
