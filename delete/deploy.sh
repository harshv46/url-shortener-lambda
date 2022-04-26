GOOS=linux GOARCH=amd64 go build -o delete main.go
zip deployment.zip delete
aws lambda create-function --function-name DeleteFunction --region ap-south-1 --zip-file fileb://./deployment.zip --runtime go1.x --tracing-config Mode=Active --role $ROLE --handler delete