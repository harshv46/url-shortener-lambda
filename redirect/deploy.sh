GOOS=linux GOARCH=amd64 go build -o redirect main.go
zip deployment.zip redirect
aws lambda create-function --function-name RedirectFunction --region ap-south-1 --zip-file fileb://./deployment.zip --runtime go1.x --tracing-config Mode=Active --role $ROLE --handler redirect
